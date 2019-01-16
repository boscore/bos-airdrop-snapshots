package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eos-go/token"
	"github.com/spf13/viper"
)

var (
	ctxKeyWorkID = 0
	ctxKeyWg     = 1
)

// Airdrop Create snapshots accounts
type Airdrop struct {
	Log *Logger

	TargetNetAPIs []*eos.API
	Snapshot      Snapshot
	SnapshotMisg  Snapshot

	Config *Config

	WriteActions bool
	ConfigFile   string
	Action       string
}

type Job struct {
	id      int
	actions []*eos.Action
}

// NewAirdrop New Airdrop process
func NewAirdrop(logger *Logger, targetAPIs []*eos.API) *Airdrop {
	a := &Airdrop{
		TargetNetAPIs: targetAPIs,
		Log:           logger,
	}
	return a
}

// CreateActions Return all the actions
func (a *Airdrop) CreateActions() (out []*eos.Action, err error) {
	snapshotFile := a.Config.Snapshot.Normal

	rawSnapshot, err := a.ReadFromCache(snapshotFile)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file: %s", err)
	}

	snapshotData, err := NewSnapshot(rawSnapshot)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot csv: %s", err)
	}

	if len(snapshotData) == 0 {
		return nil, fmt.Errorf("snapshot is empty or not loaded")
	}

	msigSnapshotData, err := NewMsigAccountSnapshot(a.Config.Snapshot.MsigJson, a.Config.MsigPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot msig json: %s", err)
	}

	if len(msigSnapshotData) == 0 {
		return nil, fmt.Errorf("msig snapshot is empty or not loaded")
	}

	a.Log.Printf("snapshot info: normal has %d rows, msig has %d rows\n", len(snapshotData), len(msigSnapshotData))

	snapshotData = append(msigSnapshotData, snapshotData...)

	for idx, hodler := range snapshotData {
		if trunc := a.Config.TestnetTruncateSnapshot; trunc != 0 && !a.Config.Mainnet {
			if idx == trunc {
				a.Log.Debugf("- DEBUG: truncated snapshot to %d rows\n", trunc)
				break
			}
		}

		creatorAccount := AN(a.Config.Creator.Name)
		destAccount := AN(hodler.BOSAccountName)
		destOwnerKey := hodler.OwnerKey
		destActiveKey := hodler.ActiveKey

		cpuStake, netStake, rest := splitSnapshotStakes(hodler.BOSBalance)
		ramQuant := uint64(3000)

		out = append(out, NewNewAccount(creatorAccount, destAccount, destOwnerKey, destActiveKey))
		out = append(out, system.NewDelegateBW(creatorAccount, destAccount, cpuStake, netStake, true))
		out = append(out, NewBuyRAM(creatorAccount, destAccount, ramQuant))
		out = append(out, token.NewTransfer(creatorAccount, destAccount, rest, "Welcome to BOS"), nil)
	}

	return
}

func splitSnapshotStakes(balance eos.Asset) (cpu, net, xfer eos.Asset) {
	// everyone has minimum 0.2000 BOS staked
	// some 10 BOS unstaked
	// the rest split between the two
	cpu = NewBOSAsset(1000)
	net = NewBOSAsset(1000)

	remainder := NewBOSAsset(balance.Amount)

	if remainder.Amount <= 100000 /* 10.0000 BOS */ {
		return cpu, net, remainder
	}

	remainder.Amount -= 100000 // keep them floating, unstaked

	firstHalf := remainder.Amount / 2
	cpu.Amount += firstHalf
	net.Amount += remainder.Amount - firstHalf

	return cpu, net, NewBOSAsset(100000)
}

// UpdateAuthActions Return all the update auth actions
func (a *Airdrop) UpdateAuthActions() (out []*eos.Action, err error) {
	// UpdateAuthActions
	msigSnapshotData, err := NewMsigSnapshot(a.Config.Snapshot.MsigJson)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot msig json: %s", err)
	}

	if len(msigSnapshotData) == 0 {
		return nil, fmt.Errorf("snapshot is empty or not loaded")
	}

	for _, hodler := range msigSnapshotData {
		destAccount := AN(hodler.BOSAccountName)
		// SORT PERMISSIONS
		var holderPerms []eos.Permission
		for _, perm := range hodler.Permissions {
			permAccts := make([]eos.PermissionLevelWeight, 0)
			permWaits := make([]eos.WaitWeight, 0)

			for _, op := range perm.RequiredAuth.Accounts {
				if op.Permission.Actor != "" {
					permAccts = append(permAccts, op)
				}
			}

			sort.Sort(PermSortByName{permAccts})

			if len(perm.RequiredAuth.Waits) != 0 {
				permWaits = perm.RequiredAuth.Waits
			}

			// Get new Authority
			permRequiredAuth := eos.Authority{
				Threshold: perm.RequiredAuth.Threshold,
				Keys:      perm.RequiredAuth.Keys,
				Accounts:  permAccts,
				Waits:     permWaits,
			}

			holderPerms = append(holderPerms, eos.Permission{
				PermName:     perm.PermName,
				Parent:       perm.Parent,
				RequiredAuth: permRequiredAuth,
			})
		}

		var ownerPerm eos.Permission
		var activePerm eos.Permission
		for _, perm := range holderPerms {
			if perm.PermName == "owner" {
				ownerPerm = perm
			} else if perm.PermName == "active" {
				activePerm = perm
			} else {
				// out = append(out, system.NewUpdateAuth(destAccount, PN(perm.PermName), PN("active"), perm.RequiredAuth, PN("owner")), nil)
			}
		}

		// Add owner and active
		out = append(out, system.NewUpdateAuth(destAccount, PN(activePerm.PermName), PN("owner"), activePerm.RequiredAuth, PN("owner")), nil)
		out = append(out, system.NewUpdateAuth(destAccount, PN(ownerPerm.PermName), PN(""), ownerPerm.RequiredAuth, PN("owner")), nil)
	}

	return
}

func (a *Airdrop) AllActions(action string) (out []*eos.Action, err error) {
	if action == "create" {
		return a.CreateActions()
	} else if action == "updateauth" {
		return a.UpdateAuthActions()
	}
	return
}

func (a *Airdrop) sendCreateTx(ctx context.Context, chks [][]*eos.Action) {
	totalChunks := len(chks)
	a.Log.Printf("Total %d chunks to send\n", totalChunks)
	jobList := make(chan Job, 1000)

	ROU := uint32(500)
	wg := new(sync.WaitGroup)
	wg.Add(int(ROU))
	for i := uint32(0); i < ROU; i++ {
		ctx := context.WithValue(ctx, ctxKeyWorkID, i)
		ctxWg := context.WithValue(ctx, ctxKeyWg, wg)
		go work(ctxWg, jobList, a.sendTx)
	}

	go func() {
		start := 0
		setp := a.Config.Tps / 2
		end := start + setp

		ticker := time.NewTicker(time.Second / 2)
		for start < totalChunks {
			time := <-ticker.C
			a.Log.Println(time.String())
			a.Log.Debugf("start = %d, end = %d, total = %d", start, end, totalChunks)

			todo := chks[start:end]
			a.Log.Printf("add %d jobs\n", len(todo))
			for i := 0; i < len(todo); i++ {
				jobList <- Job{id: i, actions: todo[i]}
				if end == totalChunks {
					a.Log.Println("All all jobs...")
				}
			}
			start += setp
			end += setp
			if end > totalChunks {
				end = totalChunks
			}
		}
		ticker.Stop()
		a.Log.Println("Push actions done")
		a.Log.Println("Waiting  for transactions to flush to blocks")

		for i := uint32(0); i < ROU; i++ {
			jobList <- Job{id: -1, actions: nil} // end work pool
		}

	}()
	wg.Wait()
}

func (a *Airdrop) sendUpdateAuthTx(chks [][]*eos.Action) error {
	for idx, chunk := range chks {
		actions := chunk
		retry := 15
		if a.Action == "updateauth" {
			retry = 0
		}
		err := Retry(retry, time.Second, func() error {
			res, err := a.SignPushActions(actions...)
			if err != nil {
				a.Log.Printf("r")
				a.Log.Debugf("error pushing transaction for chunk: %d, %s\n", idx, err)
				a.writeActions(actions, fmt.Sprintf("%s-fail.json", a.Action))
				return fmt.Errorf("push actions for chunk %d: %s", idx, err)
			} else {
				a.Log.Println("Tx: " + res.TransactionID)
				a.writeActions(actions, fmt.Sprintf("%s-success.json", a.Action))
			}
			return nil
		})
		if err != nil {
			a.Log.Println(time.Now(), " failed, try chunk...", idx, ", error: ", err)
		}
	}
	return nil
}

// Run Start airdrop process
func (a *Airdrop) Run(ctx context.Context) error {
	a.Log.Println(time.Now(), "START AIRDROP ...")
	var cache = viper.GetBool("cache")
	var acts []*eos.Action
	var err error
	var action = a.Action
	if action == "create" {
		acts, err = a.CreateActions()
	} else if action == "updateauth" {
		if cache == true {
			acts, err = GetCachedActions(action)
		} else {
			acts, err = a.UpdateAuthActions()
		}
	}

	if err != nil {
		return fmt.Errorf("getting actions: %s", err)
	}

	if err := a.writeAllActionsToDisk(); err != nil {
		return fmt.Errorf("writing %s actions to disk: %s", action, err)
	}

	if totalActs := len(acts); totalActs != 0 {
		a.Log.Printf("Total %d actions to send\n", totalActs)
		chks := ChunkifyActions(acts)
		if a.Action == "create" {
			a.sendCreateTx(ctx, chks)
		} else if a.Action == "updateauth" {
			a.sendUpdateAuthTx(chks)
		} else {
			a.Log.Println(time.Now(), "No actions to send!")
		}
		a.Log.Println(time.Now(), "All work done!")
	}

	return nil
}

func (a *Airdrop) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	if err := a.Run(ctx); err != nil {
		return err
	}
	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//<-signalChan
	cancel()
	//time.Sleep(3 * time.Second)
	return nil
}

func work(ctx context.Context, list chan Job, hook func(context.Context, Job) error) {
	fmt.Printf("%d \twork run.\n", ctx.Value(ctxKeyWorkID))
	wg, ok := ctx.Value(ctxKeyWg).(*sync.WaitGroup)
	if !ok {
		panic("ctx.Value(ctxKeyWg) error ")
	}
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%d \twork exit\n", ctx.Value(ctxKeyWorkID))
			return
		case e, ok := <-list:
			if !ok {
				fmt.Println(ctx.Value(ctxKeyWorkID), "<-job chan  fail")
			}
			if e.id < 0 {
				fmt.Printf("%d \twork exit\n", ctx.Value(ctxKeyWorkID))
				return
			}
			err := hook(ctx, e)
			if nil != err {
				fmt.Printf("%d \twork error: %v\n", ctx.Value(ctxKeyWorkID), err)
			}
		}
	}
}

func (a *Airdrop) getRandomApi() *eos.API {
	if len(a.TargetNetAPIs) == 0 {
		panic("node Api list is nil ")
	}
	nodeIdx := uint32(0)
	if len(a.TargetNetAPIs) > 1 {
		nodeIdx = rand.Uint32() % uint32(len(a.TargetNetAPIs))
	}
	return a.TargetNetAPIs[nodeIdx]
}

func (a *Airdrop) writeAllActionsToDisk() error {
	action := a.Action
	filename := fmt.Sprintf("%s-actions.json", action)
	if !a.WriteActions {
		a.Log.Printf("Not writing actions to '%s'. Activate with --write-actions\n", filename)
		return nil
	}

	a.Log.Printf("Writing all actions to '%s'...", filename)
	fl, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fl.Close()

	acts, err := a.AllActions(action)
	if err != nil {
		return fmt.Errorf("Get actions error %s", err)
	}

	for _, stepAction := range acts {
		if stepAction == nil {
			continue
		}

		stepAction.SetToServer(false)
		data, err := json.Marshal(stepAction)
		if err != nil {
			return fmt.Errorf("binary marshalling: %s", err)
		}

		_, err = fl.Write(data)
		if err != nil {
			return err
		}
		_, _ = fl.Write([]byte("\n"))
	}

	return nil
}

func (a *Airdrop) writeActions(acts []*eos.Action, filename string) error {
	fl, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer fl.Close()

	for _, stepAction := range acts {
		if stepAction == nil {
			continue
		}

		stepAction.SetToServer(false)
		data, err := json.Marshal(stepAction)
		if err != nil {
			return fmt.Errorf("binary marshalling: %s", err)
		}

		_, err = fl.Write(data)
		if err != nil {
			return err
		}
		_, _ = fl.Write([]byte("\n"))
	}
	return nil
}

func (a *Airdrop) SignPushActions(actions ...*eos.Action) (out *eos.PushTransactionFullResp, err error) {
	api := a.getRandomApi()
	opts := &eos.TxOptions{}
	if err := opts.FillFromChain(api); err != nil {
		return nil, err
	}

	tx := eos.NewTransaction(actions, opts)
	tx.SetExpiration(time.Duration(3600) * time.Second)

	return api.SignPushTransaction(tx, opts.ChainID, opts.Compress)
}

func (a *Airdrop) sendTx(ctx context.Context, job Job) error {
	idx := job.id
	actions := job.actions
	err := Retry(15, time.Second, func() error {
		res, err := a.SignPushActions(actions...)
		if err != nil {
			a.Log.Printf("r")
			a.Log.Debugf("error pushing transaction for job: %d, %s\n", idx, err)
			a.writeActions(actions, fmt.Sprintf("%s-fail.json", a.Action))
			return fmt.Errorf("push actions for job %d: %s", idx, err)
		} else {
			a.Log.Println("Tx: " + res.TransactionID)
			a.writeActions(actions, fmt.Sprintf("%s-success.json", a.Action))
		}
		return nil
	})
	if err != nil {
		a.Log.Printf(" failed\n")
		return err
	}
	return nil
}

// ChunkifyActions Chunkify actions
func ChunkifyActions(actions []*eos.Action) (out [][]*eos.Action) {
	currentChunk := []*eos.Action{}
	for _, act := range actions {
		if act == nil {
			if len(currentChunk) != 0 {
				out = append(out, currentChunk)
			}
			currentChunk = []*eos.Action{}
		} else {
			currentChunk = append(currentChunk, act)
		}
	}
	if len(currentChunk) > 0 {
		out = append(out, currentChunk)
	}
	return
}

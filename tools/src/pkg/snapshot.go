package pkg

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

type Snapshot []SnapshotLine

type SnapshotLine struct {
	EOSAccountName string
	OwnerKey       ecc.PublicKey
	ActiveKey      ecc.PublicKey
	EOSBalance     eos.Asset
	BOSAccountName string
	BOSBalance     eos.Asset
}

func NewSnapshot(content []byte) (out Snapshot, err error) {
	reader := csv.NewReader(bytes.NewBuffer(content))
	allRecords, err := reader.ReadAll()
	if err != nil {
		return
	}

	for _, el := range allRecords {
		if len(el) != 6 {
			return nil, fmt.Errorf("should have 6 elements per line")
		}

		ownerKey, err := ecc.NewPublicKey(el[1])
		if err != nil {
			return out, err
		}

		activeKey, err := ecc.NewPublicKey(el[2])
		if err != nil {
			return out, err
		}

		eosBalance, err := eos.NewEOSAssetFromString(el[3])
		if err != nil {
			return out, err
		}

		bosBalance, err := NewBOSAssetFromString(el[5])
		if err != nil {
			return out, err
		}

		out = append(out, SnapshotLine{el[0], ownerKey, activeKey, eosBalance, el[4], bosBalance})
	}

	return
}

func NewMsigAccountSnapshot(filepath, privKeyStr string) (out Snapshot, err error) {
	msigSnapshotData, err := NewMsigSnapshot(filepath)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot msig json: %s", err)
	}

	if len(msigSnapshotData) == 0 {
		return nil, fmt.Errorf("snapshot is empty or not loaded")
	}

	for _, hodler := range msigSnapshotData {

		privKey, err := ecc.NewPrivateKey(privKeyStr)
		if nil != err {
			return nil, fmt.Errorf("config privkey err: %s", privKeyStr)
		}
		pubKey := privKey.PublicKey()

		out = append(out, SnapshotLine{hodler.EOSAccountName, pubKey, pubKey, hodler.EOSBalance, hodler.BOSAccountName, hodler.BOSBalance})
	}

	return
}

type MsigSnapshot []MsigSnapshotLine
type MsigSnapshotLine struct {
	EOSAccountName string           `json:"eos_account"`
	EOSBalance     eos.Asset        `json:"eos_balance"`
	BOSAccountName string           `json:"bos_account"`
	BOSBalance     eos.Asset        `json:"bos_balance"`
	Permissions    []eos.Permission `json:"permissions"`
}

func NewMsigSnapshot(filepath string) (out MsigSnapshot, err error) {
	fi, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		var ms MsigSnapshotLine
		err = json.Unmarshal(line, &ms)
		if nil != err {
			fmt.Printf("Error: %s\n", err)
		}

		out = append(out, ms)
	}
	return
}

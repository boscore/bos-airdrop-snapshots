package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/spf13/viper"
	"gitlab.com/boscore/bos-airdrop-go/pkg"
)

func setupAirdrop(cnofigFile string, action string) (a *pkg.Airdrop, err error) {
	airdropConfig, err := pkg.ReadConfig(cnofigFile)
	if err != nil {
		return nil, err
	}

	privKeyStr := airdropConfig.Creator.PrivateKey.String()
	if action == "updateauth" {
		privKeyStr = airdropConfig.MsigPrivateKey
	}

	privKey, err := ecc.NewPrivateKey(privKeyStr)
	if nil != err {
		panic("config privkey err")
	}
	pubKey := privKey.PublicKey()

	keyBag := eos.NewKeyBag()
	if err := keyBag.Add(privKeyStr); err != nil {
		fmt.Print("Couldn't load private key:", err)
	}

	var targetNetAPIs []*eos.API
	for _, n := range airdropConfig.HttpEndpoints {
		log.Println("Start to ping api node: ", n)
		api := eos.New(n)
		if _, err := api.GetInfo(); err != nil {
			fmt.Println("init node api error : ", n, " ", err)
		} else {
			api.Signer = keyBag
			api.SetCustomGetRequiredKeys(func(tx *eos.Transaction) (keys []ecc.PublicKey, e error) {
				return []ecc.PublicKey{pubKey}, nil
			})
			api.HttpClient.Timeout = time.Duration(60) * time.Second
			api.Debug = false
			fmt.Println("In-memory keys:")
			memkeys, _ := api.Signer.AvailableKeys()
			for _, key := range memkeys {
				fmt.Printf("- %s\n", key.String())
			}
			fmt.Println("")

			targetNetAPIs = append(targetNetAPIs, api)
		}
	}

	if len(targetNetAPIs) == 0 {
		return nil, errors.New("Must have at least one `http_endpoints`.")
	}

	log.Println("Init API Node Success...")
	logger := pkg.NewLogger(action)
	logger.Debug = viper.GetBool("verbose")

	a = pkg.NewAirdrop(logger, targetNetAPIs)
	a.WriteActions = viper.GetBool("write-actions")
	a.Config = airdropConfig
	a.Action = action
	return a, nil
}

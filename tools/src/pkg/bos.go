package pkg

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
)

// EOSSymbol represents the standard EOS symbol on the chain.  It's
// here just to speed up things.
var BOSSymbol = eos.Symbol{Precision: 4, Symbol: "BOS"}

func NewBOSAsset(amount int64) eos.Asset {
	return eos.Asset{Amount: amount, Symbol: BOSSymbol}
}

// NewBOSAssetFromString return a bos asset from string
func NewBOSAssetFromString(amount string) (out eos.Asset, err error) {
	if len(amount) == 0 {
		return out, fmt.Errorf("cannot be an empty string")
	}

	if strings.Contains(amount, " BOS") {
		amount = strings.Replace(amount, " BOS", "", 1)
	}
	if !strings.Contains(amount, ".") {
		val, err := strconv.ParseInt(amount, 10, 64)
		if err != nil {
			return out, err
		}
		return eos.NewEOSAsset(val * 10000), nil
	}

	parts := strings.Split(amount, ".")
	if len(parts) != 2 {
		return out, fmt.Errorf("cannot have two . in amount")
	}

	if len(parts[1]) > 4 {
		return out, fmt.Errorf("BOS has only 4 decimals")
	}

	val, err := strconv.ParseInt(strings.Replace(amount, ".", "", 1), 10, 64)
	if err != nil {
		return out, err
	}
	return NewBOSAsset(val * int64(math.Pow10(4-len(parts[1])))), nil
}

// NewNewAccount returns a `newaccount` action that lives on the
// `eosio.system` contract.
func NewNewAccount(creator, newAccount eos.AccountName, ownerKey, activeKey ecc.PublicKey) *eos.Action {
	return &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("newaccount"),
		Authorization: []eos.PermissionLevel{
			{Actor: creator, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(system.NewAccount{
			Creator: creator,
			Name:    newAccount,
			Owner: eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{
					{
						PublicKey: ownerKey,
						Weight:    1,
					},
				},
				Accounts: []eos.PermissionLevelWeight{},
			},
			Active: eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{
					{
						PublicKey: activeKey,
						Weight:    1,
					},
				},
				Accounts: []eos.PermissionLevelWeight{},
			},
		}),
	}
}

//NewBuyRAM return a `buyram` action that lives on the
// `eosio.system` contract.
func NewBuyRAM(payer, receiver eos.AccountName, bosQuantity uint64) *eos.Action {
	a := &eos.Action{
		Account: AN("eosio"),
		Name:    ActN("buyram"),
		Authorization: []eos.PermissionLevel{
			{Actor: payer, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(system.BuyRAM{
			Payer:    payer,
			Receiver: receiver,
			Quantity: NewBOSAsset(int64(bosQuantity)),
		}),
	}
	return a
}

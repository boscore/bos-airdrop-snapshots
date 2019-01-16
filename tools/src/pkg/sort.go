package pkg

import (
	"strings"

	eos "github.com/eoscanada/eos-go"
)

type PermAccounts []eos.PermissionLevelWeight

func (p PermAccounts) Len() int      { return len(p) }
func (p PermAccounts) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type PermSortByName struct{ PermAccounts }

func (p PermSortByName) Less(i, j int) bool {
	ia := string(p.PermAccounts[i].Permission.Actor)
	ja := string(p.PermAccounts[j].Permission.Actor)

	return strings.Compare(ia, ja) < 0
}

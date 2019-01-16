package pkg

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"
)

func TestPermAccountSort(t *testing.T) {
	jsonStr := `[{"permission":{"actor":"nndfdlsheqir","permission":"active"},"weight":2},{"permission":{"actor":"fekcuoaywlye","permission":"active"},"weight":2},{"permission":{"actor":"jojzyrruhuwn","permission":"active"},"weight":2},{"permission":{"actor":"cadtcblrsteh","permission":"active"},"weight":2},{"permission":{"actor":"teirrylsgrsb","permission":"active"},"weight":1},{"permission":{"actor":"yyqzebigknvi","permission":"active"},"weight":1},{"permission":{"actor":"kddtowlpiaws","permission":"active"},"weight":1},{"permission":{"actor":"ixikjaidjtse","permission":"active"},"weight":1},{"permission":{"actor":"xwbaaujupsti","permission":"active"},"weight":1},{"permission":{"actor":"rddqajurznmb","permission":"active"},"weight":1},{"permission":{"actor":"bldxqvwuufnt","permission":"active"},"weight":1}]`
	var permAccts PermAccounts
	if err := json.Unmarshal([]byte(jsonStr), &permAccts); err == nil {
		fmt.Println("================ before sort ================")
		for _, p := range permAccts {
			fmt.Println(p.Permission.Actor, ":", p.Permission.Permission)
		}

		sort.Sort(PermSortByName{permAccts})
		fmt.Println("================ after sort ================")
		for _, p := range permAccts {
			fmt.Println(p.Permission.Actor, ":", p.Permission.Permission)
		}
	}
}

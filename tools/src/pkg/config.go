package pkg

import (
	"fmt"
	"io/ioutil"

	"github.com/eoscanada/eos-go/ecc"
)

type Config struct {
	Mainnet                 bool           `json:"mainnet"`
	TestnetTruncateSnapshot int            `json:"testnet_truncate_snapshot"`
	Tps                     int            `json:"tps"`
	HttpEndpoints           []string       `json:"http_endpoints"`
	Snapshot                ConfigSnapshot `json:"snapshot"`
	Creator                 ConfigCrtator  `json:"creator"`
	MsigPrivateKey          string         `json:"msig_prikey"`
}

type ConfigSnapshot struct {
	All      string `json:"all,omitempty"`
	Normal   string `json:"normal"`
	Msig     string `json:"msig,omitempty"`
	MsigJson string `json:"msig_json"`
}

type ConfigCrtator struct {
	Name       string          `json:"name"`
	PublicKey  ecc.PublicKey   `json:"pubkey,omitempty"`
	PrivateKey *ecc.PrivateKey `json:"prikey"`
}

func ReadConfig(filename string) (out *Config, err error) {
	rawConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config: %s", err)
	}

	if err := yamlUnmarshal(rawConfig, &out); err != nil {
		return nil, fmt.Errorf("parsing config yaml: %s", err)
	}

	return
}

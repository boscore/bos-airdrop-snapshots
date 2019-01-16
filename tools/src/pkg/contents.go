package pkg

import (
	"io/ioutil"
)

func (a *Airdrop) ReadFromCache(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

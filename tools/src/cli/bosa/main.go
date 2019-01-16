package main

import (
	// Load all contracts here, so we can always read and decode
	// transactions with those contracts.
	"gitlab.com/boscore/bos-airdrop-go/cli/bosa/cmd"
)

var version = "dev"

func init() {
	cmd.Version = version
}

func main() {
	cmd.Execute()
}

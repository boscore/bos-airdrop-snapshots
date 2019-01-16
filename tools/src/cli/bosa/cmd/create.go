package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// airdropCmd represents the API command
var airdropCmd = &cobra.Command{
	Use:   "create",
	Short: "Create accounts for BOS Mainnet.",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := "config.yaml"
		if len(args) > 0 {
			configFile = args[0]
		}

		a, err := setupAirdrop(configFile, "create")
		if err != nil {
			log.Fatalln("create setup:", err)
		}

		if err := a.Start(); err != nil {
			log.Fatalf("Create run error: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(airdropCmd)
}

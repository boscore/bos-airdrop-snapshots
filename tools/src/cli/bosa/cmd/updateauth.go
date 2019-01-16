package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// updateCmd represents the updateauth command
var updateCmd = &cobra.Command{
	Use:   "updateauth",
	Short: "Update auth for EOS Mainnet msig accounts on BOS Mainnet.",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := "config.yaml"
		if len(args) > 0 {
			configFile = args[0]
		}

		a, err := setupAirdrop(configFile, "updateauth")
		if err != nil {
			log.Fatalln("update setup:", err)
		}

		if err := a.Start(); err != nil {
			log.Fatalf("Update auth run error: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

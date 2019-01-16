// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the program version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("https://boscore.io - bosa", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

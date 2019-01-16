package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version string

// RootCmd represents the eosc command
var RootCmd = &cobra.Command{
	Use:   "bosa",
	Short: "bosa is an BOS command-line Swiss Army knife",
	Long:  `bosa is a command-line Swiss Army knife for BOS - by BOSCore.`,
}

// Execute executes the configured RootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolP("write-actions", "w", true, "Write actions to disk.")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose output (also see 'output.log')")
	RootCmd.PersistentFlags().BoolP("cache", "c", false, "If find success actions, will omit those actions.")

	for _, flag := range []string{"write-actions", "verbose", "cache"} {
		if err := viper.BindPFlag(flag, RootCmd.PersistentFlags().Lookup(flag)); err != nil {
			panic(err)
		}
	}
}

func initConfig() {
	viper.SetEnvPrefix("BOS")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
}

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Debug bool

var rootCmd = &cobra.Command{
	Use: "sheepstor",
	Run: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "--debug=true|false")
}

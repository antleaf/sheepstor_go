package cmd

import (
	. "github.com/antleaf/sheepstor/internal"
	"github.com/spf13/cobra"
)

var port int

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		InitialiseApplication(ConfigFilePath)
		RunServer(port)
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "", 8081, "--port=8081")
	rootCmd.AddCommand(serverCmd)

}

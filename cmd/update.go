package cmd

import (
	"fmt"
	. "github.com/antleaf/sheepstor/internal"
	"github.com/spf13/cobra"
)

var sites string

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		InitialiseApplication()
		Log.Info(fmt.Sprintf("Running as CLI Process, updating website(s): '%s'...", sites))
		if sites == "all" {
			Registry.ProcessAllWebsites()
		} else {
			website := *Registry.GetWebsiteByID(sites)
			website.Process()
		}
	},
}

func init() {
	updateCmd.Flags().StringVarP(&sites, "sites", "", "", "--sites all|<some_id>")
	rootCmd.AddCommand(updateCmd)
}

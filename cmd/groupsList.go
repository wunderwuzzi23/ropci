/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"ropci/utils"

	"github.com/spf13/cobra"
)

// groupsListCmd represents the groupsList command
var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		utils.DoRequest(mainClient,
			rootGraphUri,
			"", //api-version not needed for Graph API
			"groups",
			rootOutputFormat,
			rootOutputFilename,
			[]string{"id", "displayName", "description", "mailEnabled", "mail"},
			groupsSearch, //search
			rootShowAll,
			"")
	},
}

var groupsSearch string

func init() {
	groupsCmd.AddCommand(groupsListCmd)
	groupsListCmd.Flags().StringArrayVarP(&groupSelectFields, "fields", "f", []string{"id", "displayName", "description", "mailEnabled", "mail"}, "the fields to select")
	groupsListCmd.Flags().StringVarP(&groupsSearch, "search", "s", "", "Search, e.g \"displayName:ropci\" searches for all groups that contain ropci in the displayName")
	groupsListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
}

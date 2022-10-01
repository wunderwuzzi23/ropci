/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	groupsUserID      string
	groupsGroupID     string
	groupsUsername    string
	groupSelectFields []string
)

// groupsCmd represents the apps command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List or create groups",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()

	},
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}

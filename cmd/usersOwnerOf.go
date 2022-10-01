/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"ropci/utils"

	"github.com/spf13/cobra"
)

// usersOwnerOfCmd represents the usersMemberOf command
var usersOwnerOfCmd = &cobra.Command{
	Use:   "ownerof",
	Short: "List the groups and objects the given user owns ",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		getUserIsOwnerOf(usersUsername)
	},
}

var (
	usersOwnerSelectFields []string
)

func init() {
	usersCmd.AddCommand(usersOwnerOfCmd)

	usersOwnerOfCmd.Flags().StringVarP(&usersUsername, "user", "u", "", "user to get info about. If empty, logged on user.")
	usersOwnerOfCmd.Flags().StringArrayVarP(&usersOwnerSelectFields, "fields", "f", []string{"id", "@odata.type", "displayName", "mail", "description", "securityEnabled"}, "the fields to select")

}

func getUserIsOwnerOf(username string) {
	var path string
	if username == "" {
		path = "me"
	} else {
		path = "users/" + username
	}

	path += "/ownedObjects"

	utils.DoRequest(mainClient,
		rootGraphUri,
		"", //api-version not needed for Graph API
		path,
		rootOutputFormat,
		rootOutputFilename,
		usersOwnerSelectFields,
		"", //search
		rootShowAll,
		"")
}

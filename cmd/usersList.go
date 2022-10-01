/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"ropci/utils"

	"github.com/spf13/cobra"
)

var (
	usersSearch       string
	usersSelectFields []string
)

// usersListCmd represents the usersList command
var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `Note: ClientID d3590ed6-52b3-4102-aeff-aad2292ab01c can list users.`,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		utils.DoRequest(mainClient,
			rootGraphUri,
			"", //api-version not needed for Graph API
			"users",
			rootOutputFormat,
			rootOutputFilename,
			usersSelectFields,
			usersSearch, //search
			rootShowAll,
			"")
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)

	usersListCmd.Flags().StringArrayVarP(&usersSelectFields, "fields", "f", []string{"id", "userPrincipalName", "displayName", "firstName", "surname", "jobTitle", "department", "phoneNumber", "accountEnabled"}, "the fields to select")
	usersListCmd.Flags().StringVarP(&usersSearch, "search", "s", "", "Search, e.g \"displayName:ropci\" searches for all users that contain ropci in the displayName")
	usersListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
}

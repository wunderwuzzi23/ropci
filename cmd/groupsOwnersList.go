/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// groupsOwnersListCmd represents the groupsOwnersList command
var groupsOwnersListCmd = &cobra.Command{
	Use:   "owners-list",
	Short: "List owners of a group",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		path := fmt.Sprintf("/groups/%s/owners", groupsGroupID)

		utils.DoRequest(mainClient,
			rootGraphUri,
			"", //api-version not needed for Graph API
			path,
			rootOutputFormat,
			rootOutputFilename,
			[]string{"id", "displayName", "userPrincipalName", "givenName", "surname", "department", "jobTitle", "accountEnabled"},
			groupsSearch, //search
			rootShowAll,
			"")
	},
}

func init() {
	groupsCmd.AddCommand(groupsOwnersListCmd)

	groupsOwnersListCmd.Flags().StringVarP(&groupsSearch, "search", "s", "", "Search, e.g \"displayName:ropci\" searches for all groups that contain ropci in the displayName")
	groupsOwnersListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
	groupsOwnersListCmd.Flags().StringVarP(&groupsGroupID, "groupid", "g", "", "GroupID of the group to add the user to")

	groupsOwnersListCmd.MarkFlagRequired("groupid")
}

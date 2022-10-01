/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// groupsMembersListCmd represents the groupsMembersList command
var groupsMembersListCmd = &cobra.Command{
	Use:   "members-list",
	Short: "List members of a group",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		path := fmt.Sprintf("/groups/%s/members", groupsGroupID)

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
	groupsCmd.AddCommand(groupsMembersListCmd)

	groupsMembersListCmd.Flags().StringVarP(&groupsSearch, "search", "s", "", "Search, e.g \"displayName:ropci\" searches for all groups that contain ropci in the displayName")
	groupsMembersListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
	groupsMembersListCmd.Flags().StringVarP(&groupsGroupID, "groupid", "g", "", "GroupID of the group to list members")

	groupsMembersListCmd.MarkFlagRequired("groupid")
}

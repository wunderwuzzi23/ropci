/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// groupsMemberRemoveCmd represents the groupsMemberRemove command
var groupsMemberRemoveCmd = &cobra.Command{
	Use:   "members-remove",
	Short: "Remove a member from a group",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		path := "/groups/" + groupsGroupID + "/members/" + groupsUserID + "/$ref"

		utils.Delete(mainClient, rootGraphUri, path)
		fmt.Println("Done.")
	},
}

func init() {
	groupsCmd.AddCommand(groupsMemberRemoveCmd)

	groupsMemberRemoveCmd.Flags().StringVarP(&groupsUserID, "userid", "u", "", "ID of the user to be removed")
	groupsMemberRemoveCmd.Flags().StringVarP(&groupsGroupID, "groupid", "g", "", "GroupID of the relatrelevant group")
	groupsMemberRemoveCmd.MarkFlagRequired("userid")
	groupsMemberRemoveCmd.MarkFlagRequired("groupid")

}

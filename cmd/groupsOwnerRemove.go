/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// groupsOwnerRemoveCmd represents the groupsMemberRemove command
var groupsOwnerRemoveCmd = &cobra.Command{
	Use:   "owners-remove",
	Short: "Remove an owner from a group",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		path := "/groups/" + groupsGroupID + "/owners/" + groupsUserID + "/$ref"

		utils.Delete(mainClient, rootGraphUri, path)
		fmt.Println("Done.")
	},
}

func init() {
	groupsCmd.AddCommand(groupsOwnerRemoveCmd)

	groupsOwnerRemoveCmd.Flags().StringVarP(&groupsUserID, "userid", "u", "", "ID of the user to be removed")
	groupsOwnerRemoveCmd.Flags().StringVarP(&groupsGroupID, "groupid", "g", "", "GroupID of the relatrelevant group")
	groupsOwnerRemoveCmd.MarkFlagRequired("userid")
	groupsOwnerRemoveCmd.MarkFlagRequired("groupid")

}

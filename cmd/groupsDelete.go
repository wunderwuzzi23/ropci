/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// groupsDeleteCmd represents the groupsDelete command
var groupsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific group by Group ID",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		utils.Delete(mainClient, rootGraphUri, "/groups/"+groupsGroupID)
		fmt.Println("Done.")

	},
}

func init() {
	groupsCmd.AddCommand(groupsDeleteCmd)

	groupsDeleteCmd.Flags().StringVarP(&groupsGroupID, "id", "", "", "Group ID GUID of the group to delete")
	groupsDeleteCmd.MarkFlagRequired("id")

}

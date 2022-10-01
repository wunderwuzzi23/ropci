/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// groupsMemberAddCmd represents the groupsMemberAdd command
var groupsMemberAddCmd = &cobra.Command{
	Use:   "members-add",
	Short: "Add an account to a group",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		if groupsUserID != "" && groupsUsername != "" {
			fmt.Println("*** Error: Can only specify ObjectId or Username. Both were specified.")
			return
		}

		var requestBody string
		if groupsUserID != "" {
			requestBody = fmt.Sprintf(`{
			"@odata.id": "https://graph.microsoft.com/v1.0/directoryObjects/%s"
		  }`, groupsUserID)
		}

		if groupsUsername != "" {
			requestBody = fmt.Sprintf(`{
				"@odata.id": "https://graph.microsoft.com/v1.0/users/%s"
			  }`, groupsUsername)
		}

		path := fmt.Sprintf("/groups/%s/members/$ref", groupsGroupID)

		utils.Post(mainClient, rootGraphUri, path, []byte(requestBody), "application/json")

	},
}

func init() {
	groupsCmd.AddCommand(groupsMemberAddCmd)
	groupsMemberAddCmd.Flags().StringVarP(&groupsUserID, "objectid", "u", "", "id of the user/group,... to add to the group")
	groupsMemberAddCmd.Flags().StringVarP(&groupsUsername, "username", "n", "", "Try adding using upn (e.g. john@example.org")
	groupsMemberAddCmd.Flags().StringVarP(&groupsGroupID, "groupid", "g", "", "id of the relevant group")
	//groupsMemberAddCmd.MarkFlagRequired("userid")
	groupsMemberAddCmd.MarkFlagRequired("groupid")
}

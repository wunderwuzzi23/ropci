/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

var usersUserID string

// usersDeleteCmd represents the usersDelete command
var usersDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an account by UserID",
	Long:  `ClientID d3590ed6-52b3-4102-aeff-aad2292ab01c can delete users.`,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		utils.Delete(mainClient, rootGraphUri, "/users/"+usersUserID)
		fmt.Println("Done.")
	},
}

func init() {
	usersCmd.AddCommand(usersDeleteCmd)
	usersDeleteCmd.Flags().StringVarP(&usersUserID, "id", "", "", "User ID GUID of the user to delete")
	usersDeleteCmd.MarkFlagRequired("id")
}

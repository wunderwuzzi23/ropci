/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"os"
	"ropci/utils"

	"github.com/spf13/cobra"
)

var usersAddTemplate string

// usersAddCmd represents the usersAdd command
var usersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new user based on a template",
	Long:  `ClientID d3590ed6-52b3-4102-aeff-aad2292ab01c can create users.`,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		requestBody, err := os.ReadFile(usersAddTemplate)
		if err != nil {
			fmt.Println("*** template not found", err)
			return
		}

		utils.Post(mainClient, rootGraphUri, "users", requestBody, "application/json")
	},
}

func init() {
	usersCmd.AddCommand(usersAddCmd)
	usersAddCmd.Flags().StringVarP(&usersAddTemplate, "template", "t", "", "JSON file that contains details (example ./templates/usersAdd.json")
	usersAddCmd.MarkFlagRequired("template")
}

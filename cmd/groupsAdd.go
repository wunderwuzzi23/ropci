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

var groupAddTemplate string

// groupsAddCmd represents the groupsAdd command
var groupsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new group, based on the provided template. See/modify the template in the ./templates/ folder.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		requestBody, err := os.ReadFile(groupAddTemplate)
		if err != nil {
			fmt.Println("*** template not found", err)
			return
		}

		utils.Post(mainClient, rootGraphUri, "groups", requestBody, "application/json")
	},
}

func init() {
	groupsCmd.AddCommand(groupsAddCmd)

	groupsAddCmd.Flags().StringVarP(&groupAddTemplate, "template", "t", "", "JSON file that contains details (example at ./templates/groupsAdd.json")
	groupsAddCmd.MarkFlagRequired("template")
}

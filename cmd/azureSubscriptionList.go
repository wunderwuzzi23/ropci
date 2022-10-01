/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"ropci/utils"

	"github.com/spf13/cobra"
)

// subscriptionListCmd represents the subscriptionList command
var subscriptionListCmd = &cobra.Command{
	Use:   "subscriptions-list",
	Short: "List all subscriptions the account has access to",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		utils.DoRequest(azureClient,
			rootAzureMgmtUri,
			"2022-09-01",
			"subscriptions",
			rootOutputFormat,
			rootOutputFilename,
			azureSelectFields,
			"", //search
			rootShowAll,
			"")
	},
}

func init() {
	azureCmd.AddCommand(subscriptionListCmd)
	subscriptionListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
	subscriptionListCmd.Flags().StringArrayVarP(&azureSelectFields, "fields", "f", []string{"subscriptionId", "displayName", "state"}, "the fields to select")
}

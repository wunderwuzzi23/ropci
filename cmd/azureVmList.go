/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// vmListCmd represents the vmList command
var vmListCmd = &cobra.Command{
	Use:   "vm-list",
	Short: "List all VMs in a subscription",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		path := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Compute/virtualMachines",
			azureSubscriptionID)

		utils.DoRequest(azureClient,
			rootAzureMgmtUri,
			"2022-08-01",
			path,
			rootOutputFormat,
			rootOutputFilename,
			azureVmSelectFields,
			"", //search
			rootShowAll,
			"")
	},
}

func init() {
	azureCmd.AddCommand(vmListCmd)

	vmListCmd.Flags().StringVarP(&azureSubscriptionID, "subscription", "s", "", "the subscription to look in")
	vmListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
	vmListCmd.Flags().StringArrayVarP(&azureVmSelectFields, "fields", "f", []string{"id", "name", "location"}, "the fields to select")

	vmListCmd.MarkFlagRequired("subscription")
}

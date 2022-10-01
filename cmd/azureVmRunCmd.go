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

// vmRuncmdCmd represents the vmRuncmd command
var vmRuncmdCmd = &cobra.Command{
	Use:   "vm-runcmd",
	Short: "Run a command in a VM",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		path := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/runCommand?api-version=%s",
			azureSubscriptionID,
			azureResourceGroupName,
			azureVmName,
			"2022-08-01")

		azureRunCommand, err := os.ReadFile(azureCmdTemplate)
		if err != nil {
			fmt.Println("*** template not found", err)
			return
		}

		utils.Post(azureClient, rootAzureMgmtUri, path, []byte(azureRunCommand), "application/json")
		fmt.Println("Done.")

	},
}

func init() {
	azureCmd.AddCommand(vmRuncmdCmd)

	vmRuncmdCmd.Flags().StringVarP(&azureSubscriptionID, "subscription", "s", "", "the subscription to use")
	vmRuncmdCmd.Flags().StringVarP(&azureResourceGroupName, "resource-group", "r", "", "resource group to use")
	vmRuncmdCmd.Flags().StringVarP(&azureVmName, "vm-name", "n", "", "name of the VM")
	vmRuncmdCmd.Flags().StringVarP(&azureCmdTemplate, "command", "c", "", "the template file for the command to run (e.g. ./template/vmRunCmdLinux.json")

	vmRuncmdCmd.MarkFlagRequired("subscription")
	vmRuncmdCmd.MarkFlagRequired("resource-group")
	vmRuncmdCmd.MarkFlagRequired("vm-name")
	vmRuncmdCmd.MarkFlagRequired("command")

}

/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	azureClient            *http.Client
	azureSelectFields      []string
	azureSubscriptionID    string
	azureResourceGroupName string
	azureVmName            string
	azureVmSelectFields    []string
	azureCmdTemplate       string
)

// azureCmd represents the azure command
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Interact with Azure Resource Manager",
	Long: `Basic features only supported, like listing subscriptions, VMs and running commands on a VM. 
Needs an access token with scope https://management.core.windows.net//user_impersonation.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()

	},

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		readViperSettings()

		var err error
		azureClient, err = getHttpClientForClientID(
			tokenCacheAzure,
			"Azure CLI",
			"04b07795-8ddb-461a-bbee-02f9e1bf7b46",
			[]string{"openid", "offline_access", "https://management.core.windows.net//user_impersonation"})

		if err != nil {
			fmt.Println("*** Sorry, can't get access token for Azure.")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(azureCmd)
}

/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"ropci/utils"

	"github.com/spf13/cobra"
)

var appsSelectFields []string

// appsCmd represents the apps command
var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "List all the apps/clientids (servicePrincipals) available in the tenant",
	Long:  `This command can be used to retrieve clientids for testing to check if one can perform ROPC grants on them.`,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		utils.DoRequest(mainClient,
			rootGraphUri,
			"", //api-version not needed for Graph API
			"servicePrincipals",
			rootOutputFormat,
			rootOutputFilename,
			appsSelectFields,
			"", //search
			rootShowAll,
			callBody)
	},
}

// type appsCmdFlags struct {
// 	outputFilename string
// 	selectFields []string
// }

func init() {
	rootCmd.AddCommand(appsCmd)

	appsCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "writing results to this json file")
	appsCmd.Flags().StringArrayVarP(&appsSelectFields, "fields", "f", []string{"displayName", "appId", "publisherName"}, "fields to select")
}

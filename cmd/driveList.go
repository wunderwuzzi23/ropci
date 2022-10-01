/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"ropci/utils"

	"github.com/spf13/cobra"
)

var driveSelectFields []string

// driveListCmd represents the driveList command
var driveListCmd = &cobra.Command{
	Use:   "list",
	Short: "List and browse through SharePoint drives",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		if drivePath != "/" {
			drivePath = drivePath + "/children"
		}

		drivePath = "/me/drive/root/children" + drivePath

		utils.DoRequest(mainClient,
			rootGraphUri,
			"", //api-version not needed for Graph API
			drivePath,
			rootOutputFormat,
			rootOutputFilename,
			driveSelectFields,
			"", //search
			rootShowAll,
			callBody)

	},
}

func init() {
	driveCmd.AddCommand(driveListCmd)

	driveListCmd.Flags().StringVarP(&drivePath, "path", "p", "/", "path of items /")
	driveListCmd.Flags().StringArrayVarP(&driveSelectFields, "fields", "f", []string{"id", "name", "lastModifiedDateTime"}, "the fields to select, e.g. '-f name -f @microsoft.graph.downloadUrl' for download link.")
	driveListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
}

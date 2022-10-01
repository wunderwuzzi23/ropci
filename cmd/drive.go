/*
Copyright © 2022 wunderwuzzi23
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var drivePath string

// driveCmd represents the drive command
var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "List, download or upload files to Sharepoint",
	Long:  `Requires a properly scoped token, such as Microsoft Teams 1fec8e78-bce4-4aaf-ab1b-5451cc387264`,
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(driveCmd)

}

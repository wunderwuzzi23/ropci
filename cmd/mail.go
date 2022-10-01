/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	mailPath string
	mailUser string
)

// mailCmd represents the mail command
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Access mail of the user",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(mailCmd)
}

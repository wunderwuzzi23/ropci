/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// usersCmd represents the apps command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users or an individual user's details",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	usersUsername string
)

func init() {
	rootCmd.AddCommand(usersCmd)
}

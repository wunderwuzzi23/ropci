/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// authInvalidateRefreshTokenCmd represents the cmdInvalidateRefreshToken command
var authInvalidateRefreshTokenCmd = &cobra.Command{
	Use:     "invalidate",
	Aliases: []string{"logoff"},
	Short:   "Invalidate all refresh tokens for the current user",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()
		utils.Post(mainClient, rootGraphUri, "/me/invalidateAllRefreshTokens", nil, "")
		fmt.Println("Done.")
	},
}

func init() {
	authCmd.AddCommand(authInvalidateRefreshTokenCmd)
}

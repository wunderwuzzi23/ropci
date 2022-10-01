/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"ropci/utils"

	"github.com/spf13/cobra"
)

var mailSelectFields []string

// mailListCmd represents the mailList command
var mailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the mails of the user (or another user if the account has proper permissions)",
	Long:  `Add -f subject -f bodyPreview to see preview of the mail as well`,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		if mailUser == "me" {
			mailPath = "/me/messages"
		} else {
			mailPath = "/users/" + mailUser + "/messages"
		}

		utils.DoRequest(mainClient,
			rootGraphUri,
			"", //api-version not needed for Graph API
			mailPath,
			rootOutputFormat,
			rootOutputFilename,
			mailSelectFields,
			"", //search
			rootShowAll,
			"")
	},
}

func init() {
	mailCmd.AddCommand(mailListCmd)

	mailListCmd.Flags().StringVarP(&mailUser, "mail", "m", "me", "account of mailbox to access, default: me.")
	mailListCmd.Flags().StringArrayVarP(&mailSelectFields, "fields", "f", []string{"subject", "createdDateTime"}, "the fields to select")
	mailListCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
}

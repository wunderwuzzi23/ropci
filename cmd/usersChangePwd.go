/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// changepwdCmd represents the changepwd command
var changepwdCmd = &cobra.Command{
	Use:   "changepwd",
	Short: "Change the password of this user",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		if usersOldPwd == "" {
			fmt.Print("Old Password: ")
			pwdBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Print("*** Error reading password.", err)
				return
			}
			usersOldPwd = string(pwdBytes)
			fmt.Println()
		}

		if usersNewPwd == "" {
			fmt.Print("New Password: ")
			pwdBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Print("*** Error reading password.", err)
				return
			}
			usersNewPwd = string(pwdBytes)
			fmt.Println()
		}

		requestBody := []byte(fmt.Sprintf(`{"currentPassword": "%s","newPassword": "%s"}`,
			usersOldPwd,
			usersNewPwd))

		path := "me/changePassword"
		utils.Post(mainClient, rootGraphUri, path, requestBody, "application/json")
		fmt.Println("Done.")
	},
}

var (
	usersNewPwd string
	usersOldPwd string
)

func init() {
	usersCmd.AddCommand(changepwdCmd)

	changepwdCmd.Flags().StringVarP(&usersNewPwd, "new-password", "", "", "the new password of the user. If empty, will prompt")
	changepwdCmd.Flags().StringVarP(&usersOldPwd, "old-password", "", "", "the old password of the user. If empty, will prompt")
}

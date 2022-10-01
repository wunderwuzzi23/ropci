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

// setpwdCmd represents the setpwd command
var setpwdCmd = &cobra.Command{
	Use:   "setpwd",
	Short: "Set the users password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		if usersSetNewPwd == "" {
			fmt.Print("New Password: ")
			pwdBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Print("*** Error reading password.", err)
				return
			}
			usersSetNewPwd = string(pwdBytes)
			fmt.Println()
		}

		requestBody := []byte(fmt.Sprintf(`{
	"passwordProfile": {
		"forceChangePasswordNextSignIn": false,
		"password": "%s"
	}
}`, usersSetNewPwd))

		path := "users/" + usersUsername
		utils.Patch(mainClient, rootGraphUri, path, requestBody, "application/json")
		fmt.Println("Done.")
	},
}

var usersSetNewPwd string

func init() {
	usersCmd.AddCommand(setpwdCmd)
	setpwdCmd.Flags().StringVarP(&usersUsername, "user", "u", "", "the upn or ID of the user")
	setpwdCmd.Flags().StringVarP(&usersSetNewPwd, "new-password", "", "", "the new password of the user. If empty, will prompt.")

	setpwdCmd.MarkFlagRequired("user")
}

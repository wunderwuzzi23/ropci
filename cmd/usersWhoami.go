/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"encoding/json"

	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:     "who",
	Aliases: []string{"whoami", "info"},
	Short:   "Get account details of a user (no argument gives current user info)",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		var path string
		if usersUsername == "" {
			path = "me"
		} else {
			path = "users/" + usersUsername
		}

		res, err := utils.Get(mainClient, rootGraphUri, path, "application/json")
		if err != nil {
			fmt.Println("*** Error while requesting resource.")
			return
		}

		prettyJson, _ := utils.GetPrettyJSON(string(res))

		if utils.Verbose {
			fmt.Println(prettyJson)
		}

		u := &User{}
		_ = json.Unmarshal(res, &u)

		if !utils.Verbose {
			printBasicUserInfo(u)

			fmt.Println("\nUser is owner of:")
			getUserIsOwnerOf(usersUsername)

			fmt.Println("\n\nUser is member of:")
			getUserIsMemberOf(usersUsername)

			//fmt.Printf("\n*** More information on user object available. Use '-v' to show output in JSON.\n")
		}
	},
}

type User struct {
	Id                             string
	CreatedDateTime                string
	Country                        string
	Mail                           string
	OtherMails                     []string
	ProxyAddresses                 []string
	RefreshTokensValidFromDateTime string
	AccountEnabled                 bool
	BusinessPhones                 []string
	DisplayName                    string
	EmployeeId                     string
	JobTitle                       string
	MobilePhone                    string
	GivenName                      string
	Surname                        string
	UserPrincipalName              string
}

func init() {
	usersCmd.AddCommand(whoamiCmd)

	whoamiCmd.Flags().StringVarP(&usersUsername, "user", "u", "", "user to get info about. If empty, logged on user.")
}

func printBasicUserInfo(u *User) {
	fmt.Println("+-------------------------------------------------------------------------------------------------+")
	fmt.Printf("| Id:\t\t%s\n| Created:\t%s\n", u.Id, u.CreatedDateTime)
	fmt.Printf("| Upn:\t\t%s\n", u.UserPrincipalName)
	fmt.Printf("| Display Name:\t%s\n", u.DisplayName)
	fmt.Printf("| Firstname:\t%s\n| Lastname:\t%s\n", u.GivenName, u.Surname)
	fmt.Printf("| Enabled:\t%v\n", u.AccountEnabled)

	if u.EmployeeId != "" {
		fmt.Printf("| EmployeeID:\t%s\n", u.EmployeeId)
	}
	if u.JobTitle != "" {
		fmt.Printf("| Job Title:\t%v\n", u.JobTitle)
	}
	if u.Country != "" {
		fmt.Printf("| Country:\t%s\n", u.Country)
	}
	fmt.Printf("| Refresh Tkns\n| Valid Since:\t%s\n", u.RefreshTokensValidFromDateTime)
	fmt.Println("+-------------------------------------------------------------------------------------------------+")

	if u.Mail != "" {
		fmt.Printf("| Main Mail:\t%s\n", u.Mail)
	}
	if len(u.OtherMails) > 0 {
		fmt.Printf("| Other Mails:\t%v\n", u.OtherMails)
	}
	if len(u.ProxyAddresses) > 0 {
		fmt.Printf("| Proxy Addr.\t%s\n", u.ProxyAddresses)
	}
	if u.MobilePhone != "" {
		fmt.Printf("| Mobile #:\t%s\n", u.MobilePhone)
	}

	if len(u.BusinessPhones) > 0 {
		fmt.Printf("| OtherPhones: \t%v\n", u.BusinessPhones)
	}
	fmt.Println("+-------------------------------------------------------------------------------------------------+")
}

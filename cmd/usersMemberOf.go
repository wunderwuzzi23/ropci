/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

// usersMemberOfCmd represents the usersMemberOf command
var usersMemberOfCmd = &cobra.Command{
	Use:   "memberof",
	Short: "List the groups and objects the given user is a member of ",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		getUserIsMemberOf(usersUsername)
	},
}

type memberOf struct {
	Value []string
}

type getByIds struct {
	Ids   []string
	Types []string
}

var (
	usersMemberSelectFields []string
)

func init() {
	usersCmd.AddCommand(usersMemberOfCmd)

	usersMemberOfCmd.Flags().StringVarP(&usersUsername, "user", "u", "", "user to get info about. If empty, logged on user.")
	usersMemberOfCmd.Flags().StringArrayVarP(&usersMemberSelectFields, "fields", "f", []string{"id", "@odata.type", "displayName", "mail", "description"}, "the fields to select")
	usersMemberOfCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
}

func getUserIsMemberOf(username string) {
	var path string
	if username == "" {
		path = "me"
	} else {
		path = "users/" + username
	}

	path += "/getMemberObjects"

	memberBytes, err := utils.Post(mainClient, rootGraphUri, path, []byte(`{"securityEnabledOnly": false }`), "application/json")
	if err != nil {
		fmt.Println("*** Error while requesting http resource.")
		return
	}

	m := memberOf{}
	err = json.Unmarshal(memberBytes, &m)
	if err != nil {
		fmt.Println("*** Error parsing json", err)
	}

	if utils.Verbose {
		json, _ := utils.GetPrettyJSON(string(memberBytes))
		fmt.Println("Retrieved memberObject Ids:\n", json)
	}

	req := getByIds{}
	// for _, v := range m.Value {
	// 	req.Ids = append(req.Ids, v)
	// }
	req.Ids = append(req.Ids, m.Value...)

	req.Types = []string{"user", "group", "device", "application",
		"servicePrincipal", //"oauth2PermissionGrant", "appRoleAssignment",
		"administrativeUnit", "directoryRole", "orgContact"}
	requestBody, _ := json.Marshal(req)
	if utils.Verbose {
		fmt.Println("Requesting details:\n\n", string(requestBody))
	}

	utils.DoRequest(mainClient,
		rootGraphUri,
		"", //api-version not needed for Graph API
		"directoryObjects/getByIds",
		rootOutputFormat,
		rootOutputFilename,
		usersMemberSelectFields,
		"", //search
		rootShowAll,
		string(requestBody))
}

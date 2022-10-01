/*
Copyright © 2022 wunderwuzzi23
*/
package cmd

import (
	"fmt"
	"ropci/utils"

	"github.com/spf13/cobra"
)

var (
	callPath         string
	callVerb         string
	callBody         string
	callSearch       string
	callSelectFields []string
)

// callCmd represents the call command
var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Generic call method to invoke arbitrary APIs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		if callVerb == "GET" {

			utils.DoRequest(mainClient,
				rootGraphUri,
				"", //api-version not needed for Graph API
				callPath,
				rootOutputFormat,
				rootOutputFilename,
				callSelectFields,
				callSearch, //search
				rootShowAll,
				"")

		} else {

			//Prefer: outlook.body-content-type="text"
			res, _ := utils.GenericRequest(mainClient, callVerb, rootGraphUri+"/"+callPath, []byte(callBody), "application/json")
			if res != nil {
				fmt.Println(string(res))
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(callCmd)

	callCmd.Flags().StringVarP(&callPath, "command", "c", "/me", "query/command to invoke")
	callCmd.Flags().StringVarP(&callSearch, "search", "s", "", "Search, e.g \"displayName:ropci\" searches for all groups that contain ropci in the displayName")
	callCmd.Flags().StringVarP(&callVerb, "verb", "", "GET", "GET, POST")
	callCmd.Flags().StringVarP(&callBody, "body", "b", "", "request body")
	callCmd.Flags().StringArrayVarP(&callSelectFields, "fields", "f", []string{"id", "displayName", "name"}, "the fields to select, e.g. -f id -f surname")
	callCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
}

/*
Copyright © 2022 wunderwuzzi23
*/
package cmd

import (
	"fmt"
	"ropci/utils"
	"strings"

	"github.com/spf13/cobra"
)

var (
	searchType         []string
	searchQueryString  string
	searchSelectFields []string
	searchMaxResults   int
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&searchQueryString, "query", "q", "password", "string to search for")
	searchCmd.Flags().StringArrayVarP(&searchType, "type", "t", []string{"message"}, "One or more of: list, site, listItem, message, event, drive, driveItem, person, externalItem")
	searchCmd.Flags().StringArrayVarP(&searchSelectFields, "fields", "f", []string{"rank", "summary", "subject"}, "the fields to select")
	searchCmd.Flags().IntVarP(&searchMaxResults, "max", "", 100, "the number of hits to show")
	searchCmd.Flags().StringVarP(&rootOutputFilename, "output", "o", "", "write the results to this file (based on --format)")
}

// searchCmd represents the drive command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search through mail messages, chats or Sharepoint files",
	Long:  `Requires a properly scoped token, such as Office 365 Management 00b41c95-dab0-4487-9791-b9d2c32c80f2`,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()

		//requestBody := fmt.Sprintf(`{"requests":[{"entityTypes":["message"],"query":{"queryString":"%s"},"from":0,"size":%d}]}`, searchQueryString, 100)

		requestBody := fmt.Sprintf(`{"requests":[{"entityTypes":["%s"],"query":{"queryString":"%s"},"from":0,"size":%d}]}`,
			strings.Join(searchType, "\",\""), searchQueryString, searchMaxResults)

		utils.DoRequest(mainClient,
			rootGraphUri,
			"", //api-version not needed for Graph API
			"search/query",
			rootOutputFormat,
			rootOutputFilename,
			searchSelectFields,
			"", //search for fields with certain values (not search API that is called here)
			rootShowAll,
			requestBody)

	},
}

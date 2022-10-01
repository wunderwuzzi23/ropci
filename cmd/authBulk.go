/*
Copyright © 2022 wunderwuzzi23
*/
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"ropci/models"
	"ropci/utils"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	bulkInputfile  string
	bulkOutputfile string
)

// authBulkCmd represents the bulk command
var authBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk validation of access to clientids via ROPC",
	Long:  `Allows to highlight OAuth apps that might allow elevation or MFA bypass opportunities.`,
	Run: func(cmd *cobra.Command, args []string) {

		processViper()

		// inputfile = viper.GetString("inputfile")
		// outputfile = viper.GetString("outputfile")

		fmt.Printf("ClientIDs from CSV file %s.\n", bulkInputfile)
		fmt.Printf("Results will be written to %s.\n\n", bulkOutputfile)

		validateClientIDs(bulkInputfile, bulkOutputfile)
	},
}

func init() {
	authCmd.AddCommand(authBulkCmd)

	authBulkCmd.Flags().StringVarP(&bulkInputfile, "inputfile", "i", "clients.csv", "CSV file with Displayname,ClientID, ClientSecret(optional)")
	authBulkCmd.Flags().StringVarP(&bulkOutputfile, "outputfile", "o", "results.json", "File with the authentication results)")

	authBulkCmd.MarkFlagRequired("inputfile")
	authBulkCmd.MarkFlagRequired("outputfile")
}

func validateClientIDs(clientidsCSVFile string, outputfile string) {

	//Create Result file
	outfile, err := os.OpenFile(outputfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer outfile.Close()

	tokch := make(chan authLogonResult, 20) // create a channel

	//simple wait group
	count := 0

	file, err := os.Open(clientidsCSVFile)
	if err != nil {
		fmt.Println("*** Error opening file", err)
	}
	defer file.Close()

	csvLines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for _, line := range csvLines {

		// for _, line := range utils.ReadFileAsStringArray(clientidsCSVFile) {
		//comp := strings.Split(line, ",")
		if len(line) < 2 {
			fmt.Println("*** Error: Incorrect csv format. Need DisplayName,ClientID.")
			return
		}

		authDisplayName := line[0]
		authClientID := line[1]
		authClientSecret := ""
		// if len(line) > 2 { //secret is most likely not included (also not needed)
		// 	authClientSecret = line[2]
		// }

		conf := &models.OAuth2Config{
			Username:      authUsername,
			Password:      authPassword,
			ClientID:      authClientID,
			ClientSecret:  authClientSecret,
			Scopes:        authScopes,
			TokenEndpoint: utils.GetAzureTokenEndpoint(authTenant),
		}

		//go ropcAuthenticate(tokch, conf, "")
		go logon(tokch, conf, authDisplayName)

		//go ropcAuthenticate(tokch, conf, authDisplayName)
		//go logon(tokch, conf, authDisplayName)

		count++

		if count%20 == 0 {
			fmt.Printf("\033[2K\rIssuing Requests...~%d", count)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{"displayName", "appId", "result", "scope"})

	fmt.Println("\nWaiting for results...")

	//collect results
	for i := 0; i < count; i++ {
		result := <-tokch
		t := result.Token

		//Authenticate
		//t, err := ropcAuthenticate(conf, displayName)
		if t.Error != "" {

			table.Append([]string{t.DisplayName, t.ClientID, "error", ""})
			//fmt.Printf("%s, %s, %s\n", t.DisplayName, t.ClientID, "Error")
			//t.Error = err.Error()

			line, err := json.MarshalIndent(t, "", "  ")
			if err != nil {
				fmt.Println("*** Error: Marshal ", err)
			}

			c, err := outfile.WriteString(string(line))
			if err != nil {
				fmt.Printf("*** Error: outfile %v,%v\n\n", c, err)
			}

			continue
		}

		s := strings.Replace(t.Scope, "00000003-0000-0000-c000-000000000000/", "", -1)
		// s = strings.Replace(s, " ", "\n", -1)

		if utils.Verbose {
			table.Append([]string{t.DisplayName, t.ClientID, "success", s, t.AccessToken})
		} else {
			table.Append([]string{t.DisplayName, t.ClientID, "success", s})
		}

		line, _ := json.Marshal(t)
		outfile.WriteString(string(line))
	}

	table.Render()

	fmt.Println("\nDone. ")
	fmt.Printf(`You could now run the following command to analyze valid tokens and there scopes:
$ cat %s | jq -r 'select (.access_token!="") | [.display_name,.scope] | @csv'`, outputfile)
	fmt.Println("\nHappy Hacking.")

}

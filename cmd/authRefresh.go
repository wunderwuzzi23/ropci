/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ropci/models"
	"ropci/utils"
	"time"

	"github.com/spf13/cobra"
)

var (
	authRefreshToken string
	authResource     string
)

// authRefreshCmd represents the refresh command
var authRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Get a new access token using the refresh token",
	Long:  `By default the refresh token in .token file is used, but it can be overwritten using --token`,
	Run: func(cmd *cobra.Command, args []string) {

		readViperSettings()

		if authRefreshToken == "" {

			t, err := readTokenFromFile(rootTokenCacheFile)
			if err != nil {
				fmt.Println("No refresh token found as argument or from tokenfile", err)
				return
			}

			if t.RefreshToken == "" {
				fmt.Println("No refresh token found as argument or from tokenfile")
				return
			}

			authRefreshToken = t.RefreshToken
		}

		fmt.Println("Attempting refresh using refresh_token")

		refresh()
	},
}

func init() {
	authCmd.AddCommand(authRefreshCmd)
	authRefreshCmd.Flags().StringVarP(&authRefreshToken, "refresh-token", "r", "", "use refresh token to get new access token")
	authRefreshCmd.Flags().StringVarP(&authResource, "resource", "", "", "add resource parameter")
}

func refresh() {

	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", authRefreshToken)
	form.Add("client_id", authClientID)
	if authResource != "" {
		form.Add("resource", authResource)
	}
	if authClientSecret != "" {
		form.Add("client_secret", authClientSecret)
	}

	if utils.Verbose {
		fmt.Println(form)
	}

	client := &http.Client{
		Transport: &utils.UserAgentRoundTripper{},
	}
	req, err := client.PostForm(utils.GetAzureTokenEndpoint(authTenant), form)
	if err != nil {
		fmt.Printf("*** Error: %v\n", err)
		return
	}

	req.Header.Set("User-Agent", utils.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencode")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("*** Error: %v\n", err)
		fmt.Printf("%s\n", body)
		return
	}

	if req.StatusCode != 200 {
		fmt.Printf("*** StatusCode: %v\n", req.Status)
		fmt.Printf("%s\n", body)
		fmt.Println("*** No new token retrieved.")
		return
	}

	t := models.Token{Time: time.Now().Format("2006-01-02 15:04:05"), ClientID: authClientID, DisplayName: authDisplayName}
	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Printf("*** Error occured unmarshal: %v \n", err)
		return
	}

	seconds := time.Duration(t.ExpiresIn - 60*5)
	t.Expiry = time.Now().Add(time.Second * seconds)

	utils.PersistTokenToFile(&t, rootTokenCacheFile)
}

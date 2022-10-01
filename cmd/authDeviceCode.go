/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ropci/utils"
	"time"

	"ropci/models"

	"github.com/spf13/cobra"
)

// authDeviceCodeCmd represents the authDeviceCode command
var authDeviceCodeCmd = &cobra.Command{
	Use:   "devicecode",
	Short: "Authenticate using device code flow (non ROPC)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		readViperSettings()

		dc, err := authUsingDeviceCodeInitiate()
		if err != nil {
			fmt.Println("*** Stopping auth flow.")
			return
		}

		fmt.Println(dc.Message)

		var token *models.Token
		for i := 0; i < 20; i++ {

			time.Sleep(time.Duration(dc.Interval) * time.Second)

			token, err = authUsingDeviceCodeGetAccessToken(dc.DeviceCode)
			if err != nil {
				fmt.Print(".")
				continue
			}

			break
		}

		if token != nil {

			seconds := time.Duration(token.ExpiresIn - 60*5)
			token.Expiry = time.Now().Add(time.Second * seconds)

			utils.PersistTokenToFile(token, rootTokenCacheFile)
		} else {
			fmt.Println("*** Sorry, device code flow did not succeed.")
			fmt.Println("*** No new access token retrieved.")
		}
	},
}

func init() {
	authCmd.AddCommand(authDeviceCodeCmd)
}

func authUsingDeviceCodeInitiate() (*models.DeviceCodeResponse, error) {

	form := url.Values{}
	form.Add("client_id", authClientID)

	scopes := ""
	for _, s := range authScopes {
		scopes += s + " "
	}
	form.Add("scope", scopes)

	if authClientSecret != "" {
		form.Add("client_secret", authClientSecret)
	}

	if utils.Verbose {
		fmt.Println(form)
	}

	client := &http.Client{
		Transport: &utils.UserAgentRoundTripper{},
	}

	req, err := client.PostForm(utils.GetAzureDeviceCodeEndpoint(authTenant), form)
	if err != nil {
		fmt.Printf("*** Error: %v\n", err)
		return nil, err
	}

	req.Header.Set("User-Agent", utils.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencode")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("*** Error: %v\n", err)
		fmt.Printf("%s\n", body)
		return nil, err
	}

	if req.StatusCode != 200 {

		fmt.Printf("*** StatusCode: %v\n", req.Status)
		fmt.Printf("%s\n", body)

		return nil, errors.New(req.Status)
	}

	deviceCode := &models.DeviceCodeResponse{}
	err = json.Unmarshal(body, &deviceCode)
	if err != nil {
		fmt.Printf("*** Error occured unmarshal: %v \n", err)
		return nil, err
	}

	return deviceCode, nil
}

func authUsingDeviceCodeGetAccessToken(deviceCode string) (*models.Token, error) {

	form := url.Values{}
	form.Add("grant_type", "device_code")
	form.Add("client_id", authClientID)

	scopes := ""
	for _, s := range authScopes {
		scopes += s + " "
	}
	form.Add("scope", scopes)

	if authClientSecret != "" {
		form.Add("client_secret", authClientSecret)
	}

	form.Add("code", deviceCode)

	if utils.Verbose {
		fmt.Println(form)
	}

	client := &http.Client{
		Transport: &utils.UserAgentRoundTripper{},
	}

	req, err := client.PostForm(utils.GetAzureTokenEndpoint(authTenant), form)
	if err != nil {
		fmt.Printf("*** Error: %v\n", err)
		return nil, err
	}

	req.Header.Set("User-Agent", utils.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencode")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("*** Error: %v\n", err)
		fmt.Printf("%s\n", body)
		return nil, err
	}

	if req.StatusCode != 200 {
		if utils.Verbose {
			fmt.Printf("*** StatusCode: %v\n", req.Status)
			fmt.Printf("%s\n", body)
		}
		return nil, errors.New(req.Status)
	}

	t := &models.Token{Time: time.Now().Format("2006-01-02 15:04:05"), ClientID: authClientID, DisplayName: authDisplayName}

	err = json.Unmarshal(body, &t)
	if err != nil {
		fmt.Printf("*** Error occured unmarshal: %v \n", err)
		return nil, err
	}

	return t, nil
}

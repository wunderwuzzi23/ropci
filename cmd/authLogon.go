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
	"time"

	"ropci/utils"

	"github.com/spf13/cobra"
)

// logonCmd represents the logon command
var logonCmd = &cobra.Command{
	Use:   "logon",
	Short: "Retrieve and store an access_token for later use.",
	Long:  `This will authenticate to the clientid and store the retrieved token information in .token`,
	Run: func(cmd *cobra.Command, args []string) {
		readViperSettings()
		authenticate()
	},
}

var authLogonDiscardToken bool

func init() {
	authCmd.AddCommand(logonCmd)

	logonCmd.Flags().BoolVarP(&authLogonDiscardToken, "discard-token", "D", false, "If authentication succeeds, discard access token and do not store it in .token file.")
}

func authenticate() {

	conf := &models.OAuth2Config{
		Username:      authUsername,
		Password:      authPassword,
		ClientID:      authClientID,
		ClientSecret:  authClientSecret,
		Scopes:        authScopes,
		TokenEndpoint: utils.GetAzureTokenEndpoint(authTenant),
	}

	tokch := make(chan authLogonResult, 1)

	//go ropcAuthenticate(tokch, conf, "")

	go logon(tokch, conf, "Main Token")

	result := <-tokch
	t := result.Token

	if t.Error != "" {
		fmt.Printf("*** Error (ClientID: %s): %s\n", t.ClientID, t.Error)
		fmt.Println("*** No new token retrieved.")
		return
	}

	if !authLogonDiscardToken {
		utils.PersistTokenToFile(t, rootTokenCacheFile)
	} else {
		fmt.Printf("Succeeded. Access token retrieved for ClientID: %s.\n", authClientID)
		fmt.Printf("Token will not be persisted.\n\n")
		fmt.Printf("Retrieved scopes:\n%s\n", t.Scope)
	}
}

// consider switching back to built in oauth2 library call
func logon(tokch chan authLogonResult, conf *models.OAuth2Config, displayName string) {

	t := &models.Token{Time: time.Now().Format("2006-01-02 15:04:05"), ClientID: conf.ClientID, DisplayName: displayName}

	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", conf.Username)
	form.Add("password", conf.Password)
	form.Add("client_id", conf.ClientID)

	scopes := ""
	for _, s := range conf.Scopes {
		scopes += s + " "
	}
	form.Add("scope", scopes)

	if conf.ClientSecret != "" {
		form.Add("client_secret", conf.ClientSecret)
	}

	// comment, since this prints the password also
	// if utils.Verbose {
	// 	fmt.Println(form)
	// }

	r := authLogonResult{Token: t, Conf: conf}

	client := &http.Client{
		Transport: &utils.UserAgentRoundTripper{},
	}

	// if utils.Verbose {
	// 	fmt.Printf("\nForm: %s\nTokenEndpoint: %s\n", form, conf.TokenEndpoint)
	// }

	req, err := client.PostForm(conf.TokenEndpoint, form)
	if err != nil {
		if utils.Verbose {
			fmt.Println(err)
		}
		r.Token.Error = err.Error()
		tokch <- r
		return
	}

	req.Header.Set("User-Agent", utils.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencode")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		if utils.Verbose {
			fmt.Println(err)
		}
		r.Token.Error = err.Error()
		tokch <- r
		return
	}

	if req.StatusCode != 200 {

		if utils.Verbose {
			fmt.Println(req.Status)
			fmt.Println(string(body))
		}
		r.Token.Error = req.Status + " " + string(body)
		tokch <- r
		return
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		if utils.Verbose {
			fmt.Println(err)
		}
		r.Token.Error = fmt.Sprintf("*** Error occured unmarshal: %v \n", err)
		tokch <- r
		return
	}

	seconds := time.Duration(t.ExpiresIn - 60*5)
	r.Token.Expiry = time.Now().Add(time.Second * seconds)

	tokch <- r
}

// func ropcAuthenticate(tokch chan models.Token, conf *oauth2.Config, displayName string) {

// 	//ctx := context.Background()

// 	t := models.Token{Time: time.Now().Format("2006-01-02 15:04:05"), ClientID: conf.ClientID, DisplayName: displayName}

// 	//fmt.Printf("Username: %s, Password: %s\n", username, password)
// 	//Authenticate
// 	token, err := conf.PasswordCredentialsToken(ctx, authUsername, authPassword)
// 	if err != nil {
// 		t.Error = err.Error()
// 		tokch <- t
// 		return
// 	}

// 	t.Scope = fmt.Sprintf("%v", token.Extra("scope"))
// 	t.AccessToken = token.AccessToken
// 	t.RefreshToken = token.RefreshToken
// 	t.IDToken = fmt.Sprintf("%v", token.Extra("id_token"))
// 	foci := token.Extra("foci")
// 	if foci == nil {
// 		foci = ""
// 	}
// 	t.Foci = fmt.Sprintf("%s", foci)
// 	t.ExpiresIn, _ = strconv.Atoi(fmt.Sprintf("%v", token.Extra("expires_in")))

// 	tokch <- t

// }

/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"ropci/models"
	"ropci/utils"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/term"
)

var (
	rootVerbose        bool
	rootGraphUri       string
	rootAzureMgmtUri   string
	rootShowAll        bool
	rootShowVersion    bool
	rootOutputFormat   string
	rootOutputFilename string

	mainClient         *http.Client
	rootTokenCacheFile string

	cfgFile     string
	VersionInfo string
)

const (
	tokenCacheFile  = ".token"
	tokenCacheMail  = ".token-mailsend"
	tokenCacheAzure = ".token-azure"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ropci",
	Short: "Resource Owner Password Credentials Assessment Tool for AAD.",
	Long: `Resource Owner Password Credentials Assessment Tool for AAD.
ropci by wunderwuzzi23`,

	Run: func(cmd *cobra.Command, args []string) {

		if rootShowVersion {
			fmt.Printf("%s [Happy hacking]\n", cmd.Version)
			os.Exit(0)
		}

		cmd.Help()
	},
	Version: VersionInfo,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is ./.ropci.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&rootVerbose, "verbose", "v", false, "Print more info when running commands")
	rootCmd.PersistentFlags().BoolVarP(&rootShowAll, "all", "a", false, "Retrieve all records for 'list' sub-commands, this could take a while")
	rootCmd.PersistentFlags().StringVarP(&rootOutputFormat, "format", "", "table", "Output results in table, csv or json (when applicable)")
	rootCmd.PersistentFlags().StringVarP(&rootGraphUri, "graphuri", "G", "https://graph.microsoft.com/beta", "Graph API endpoint/version to call")
	rootCmd.PersistentFlags().StringVarP(&rootAzureMgmtUri, "azureuri", "A", "https://management.azure.com", "Azure Resource Management Uri")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// Search config in home directory with name ".ropci" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ropci")
	}

	// Disabling environment variable processing.
	// It caused issues on Windows as username is a default environment variable
	// viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {

		if utils.Verbose {
			fmt.Println(err.Error())
		}

		//configureRopci()
	}
}

func readTokenFromFile(filename string) (*models.Token, error) {
	var token *models.Token

	tokenfile, err := os.ReadFile(filename)
	if err != nil {
		if utils.Verbose {
			fmt.Printf("*** Cached token not found: %s\n", filename)
		}
		return nil, err
	}

	if utils.Verbose {
		fmt.Printf("*** Found %s file.\n", filename)
	}

	err = json.Unmarshal([]byte(tokenfile), &token)
	if err != nil {
		if utils.Verbose {
			fmt.Printf("*** Error: Not authenticated. No valid token in token file.\n %v\n", err)
		}

		return nil, err
	}

	if token != nil && token.AccessToken == "" {
		if utils.Verbose {
			fmt.Printf("*** Error: No access token found in token file.\n")
		}

		return nil, err
	}

	return token, nil
}

func processViper() {

	readViperSettings()

	var err error
	mainClient, err = getHttpClientForClientID(rootTokenCacheFile, authDisplayName, authClientID, authScopes)
	if err != nil {
		os.Exit(1)
	}

	//readViperSettings()
}

func readViperSettings() {

	//The assignment via viper for these didn't work automaticlly in case
	//the flag is set in the config file - so doing this manually.
	authTenant = viper.GetString("tenant")
	authUsername = viper.GetString("username")
	authPassword = viper.GetString("password")
	authClientID = viper.GetString("clientid")
	authClientSecret = viper.GetString("clientsecret")
	authScopes = viper.GetStringSlice("scope")
	authEnterPassword = viper.GetBool("enter-password")
	rootGraphUri = viper.GetString("graphuri")
	rootAzureMgmtUri = viper.GetString("azureuri")

	utils.Verbose = rootVerbose
	utils.UserAgent = viper.GetString("useragent")

	rootTokenCacheFile = tokenCacheFile

	if authTenant == "" {
		fmt.Printf("No tenant information provided in config file or on command line.\n\n")
		fmt.Printf("It looks like you don't have a valid configuration setup yet.\n\nTo configure ropci run:\n\t./ropci configure\n\n")
		fmt.Printf("To perform a quick ROPC test for an account run:\n\t./ropci auth logon -t {tenant}.onmicrosoft.com -u {user@domain.org} -P --discard-token\n\n")
		fmt.Printf("Happy hacking.\n")
		os.Exit(1)
	}

	if authUsername == "" {
		fmt.Println("No username information provided in config file or on command line.")
		os.Exit(1)
	}

	if utils.Verbose {
		fmt.Printf("Tenant:    %s\n", authTenant)
		fmt.Printf("Username:  %s\n", authUsername)
		fmt.Printf("ClientID:  %s\n", authClientID)
		fmt.Printf("Scopes:    %s\n", authScopes)
	}

	if authEnterPassword || authPassword == "" {
		fmt.Print("Password: ")
		pwdBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Print("*** Error reading password.", err)
			os.Exit(1)
		}
		authPassword = string(pwdBytes)
		fmt.Println()
	}

}

func getHttpClientForClientID(tokenCacheFilename string, displayName string, clientID string, scope []string) (*http.Client, error) {
	cachedToken, err := readTokenFromFile(tokenCacheFilename)
	if err != nil || cachedToken == nil {
		if utils.Verbose {
			fmt.Printf("*** Not authenticated yet for %s.\n", displayName)
		}
	}

	//set oauth2 config up to leverage built in refresh capabilities
	c := &oauth2.Config{
		//Scopes: []string{"openid", "offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  utils.GetAzureAuthEndpoint(authTenant),
			TokenURL: utils.GetAzureTokenEndpoint(authTenant),
		},
	}

	hc := &http.Client{Transport: &utils.UserAgentRoundTripper{}}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, hc)
	tk := &oauth2.Token{TokenType: "Bearer"}

	if cachedToken == nil {
		fmt.Printf("Attempting to get an access token.")

		conf := &models.OAuth2Config{
			Username:      authUsername,
			Password:      authPassword,
			ClientID:      clientID,
			Scopes:        scope,
			TokenEndpoint: utils.GetAzureTokenEndpoint(authTenant),
		}

		tokch := make(chan authLogonResult, 1) // create a channel

		go logon(tokch, conf, displayName)

		res := <-tokch
		t := res.Token

		if t.Error != "" {
			//fmt.Printf("*** Error (ClientID: %s): %s\n", t.ClientID, t.Error)
			return nil, errors.New(t.Error)
		}

		utils.PersistTokenToFile(t, tokenCacheFilename)

		tk.AccessToken = t.AccessToken
		tk.RefreshToken = t.RefreshToken
		tk.Expiry = t.Expiry

		return c.Client(ctx, tk), nil
	}

	if utils.Verbose {
		fmt.Println("*** Found cached token. Checking if it is expired...")
	}

	restoredToken := &oauth2.Token{
		AccessToken:  cachedToken.AccessToken,
		RefreshToken: cachedToken.RefreshToken,
		Expiry:       cachedToken.Expiry,
		TokenType:    "Bearer",
	}

	ts := c.TokenSource(context.Background(), restoredToken)
	_ = oauth2.NewClient(context.Background(), ts)
	tk, err = ts.Token()
	if err != nil {

		if utils.Verbose {
			fmt.Println("*** unexpected refresh error ", err)
		}
		return nil, err
	}

	//check if token got refreshed?
	if cachedToken.AccessToken != tk.AccessToken {

		if utils.Verbose {
			fmt.Println("*** New access token retrieved.")
		}

		//update cachedToken and store the new one
		cachedToken.AccessToken = tk.AccessToken
		cachedToken.RefreshToken = tk.RefreshToken
		cachedToken.IDToken = fmt.Sprintf("%v", tk.Extra("id_token"))
		cachedToken.Scope = fmt.Sprintf("%v", tk.Extra("scope"))

		foci := tk.Extra("foci")
		if foci == nil {
			foci = ""
		}
		cachedToken.Foci = fmt.Sprintf("%s", foci)
		cachedToken.ExpiresIn, _ = strconv.Atoi(fmt.Sprintf("%v", tk.Extra("expires_in")))

		seconds := time.Duration(cachedToken.ExpiresIn - 60*5)
		cachedToken.Expiry = time.Now().Add(time.Second * seconds)

		fmt.Printf("Access token was refreshed\n")
		utils.PersistTokenToFile(cachedToken, tokenCacheFilename)

	} else {

		if utils.Verbose {
			fmt.Println("*** Not expired. Using cached token.")
		}

		tk.AccessToken = cachedToken.AccessToken
		tk.RefreshToken = cachedToken.RefreshToken
		tk.Expiry = cachedToken.Expiry
	}

	return c.Client(ctx, tk), nil
}

/*
Copyright © 2022 wunderwuzzi23
*/
package cmd

import (
	"ropci/models"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	authTenant        string
	authClientID      string
	authClientSecret  string
	authUsername      string
	authPassword      string
	authScopes        []string
	authDisplayName   string
	authEnterPassword bool
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate to AAD using ROPC",
	Long:  `Provide username, password and clientid.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.PersistentFlags().StringP("tenant", "t", "", "AAD Tenant name or ID")
	authCmd.PersistentFlags().StringP("username", "u", "", "Username for tenant")
	authCmd.PersistentFlags().StringP("password", "p", "", "Password for user")
	authCmd.PersistentFlags().StringP("clientid", "c", "d3590ed6-52b3-4102-aeff-aad2292ab01c", "ClientID")
	authCmd.PersistentFlags().StringP("clientsecret", "s", "", "optional (not needed for basic tests)")
	authCmd.PersistentFlags().StringSliceP("scope", "S", []string{"openid", "offline_access"}, "Requested scopes")
	authCmd.PersistentFlags().BoolP("enter-password", "P", false, "Prompt for password")

	//authCmd.MarkPersistentFlagRequired("tenant")
	viper.BindPFlags(authCmd.PersistentFlags())

	authCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if viper.IsSet(flag.Name) && viper.GetString(flag.Name) != "" {
			authCmd.PersistentFlags().Set(flag.Name, viper.GetString(flag.Name))
			authCmd.PersistentFlags().SetAnnotation(flag.Name, cobra.BashCompOneRequiredFlag, []string{"false"})
		}
	})

}

type authLogonResult struct {
	wg    *sync.WaitGroup
	Token *models.Token
	Conf  *models.OAuth2Config
}

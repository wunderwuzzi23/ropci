/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"ropci/utils"

	"github.com/spf13/cobra"
)

var mailMessageTemplate string

// mailSendCmd represents the mailSend command
var mailSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a mail as the user (or another user) using a pre-defined mail template.",
	Long:  `Needs proper clientid with correct scopes, for sending mail 57336123-6e14-4acc-8dcf-287b6088aa28`,
	Run: func(cmd *cobra.Command, args []string) {

		readViperSettings()

		if mailUser == "me" {
			mailPath = "/me/sendMail"
		} else {
			mailPath = "/users/" + mailUser + "/sendMail"
		}

		mailClient, err := getHttpClientForClientID(
			tokenCacheMail,
			"Mail.Send (Microsoft Whiteboard Client)",
			"57336123-6e14-4acc-8dcf-287b6088aa28",
			[]string{"openid", "offline_access"})

		if err != nil {
			fmt.Println("*** Sorry, can't send mail.")
			return
		}

		send(mailClient)

	},
}

func init() {
	mailCmd.AddCommand(mailSendCmd)

	mailSendCmd.Flags().StringVarP(&mailUser, "mail", "m", "me", "account of mailbox to access, default: me.")
	mailSendCmd.Flags().StringVarP(&mailMessageTemplate, "template", "t", "", "JSON file that contains the message details (example at ./templates/sendMail.json")
	mailSendCmd.MarkFlagRequired("template")
}

// TODO: switch to utils.Post
func send(client *http.Client) {

	fullUri := rootGraphUri + mailPath

	requestBody, err := os.ReadFile(mailMessageTemplate)
	if err != nil {
		fmt.Println("*** mail template not found", err)
	}

	if utils.Verbose {
		fmt.Printf("Request:  %s\n", fullUri)

		if len(requestBody) > 0 {
			fmt.Printf("POST Body: %v\n", string(requestBody))
		}
	}

	var resp *http.Response
	var req *http.Request

	req, err = http.NewRequest("POST", fullUri, bytes.NewBufferString(string(requestBody)))
	if err != nil {
		fmt.Println("*** Error creating request", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", utils.UserAgent)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("*** Error issuing request: ", err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("*** Error occured while reading content: " + err.Error())
	}

	if resp.StatusCode != 202 {
		fmt.Println(string(body))
		return
	}

	fmt.Println("Mail sent.")
}

/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"ropci/utils"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// configureCmd represents the config command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Initialize the ropci configuration to get started.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		configureRopci()

	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func configureRopci() {

	//not using viper.WriteConfig because it does not
	//maintain the order and layout of the yaml file
	fmt.Println("Let's set things up by entering Tenant, Username and Password.")
	//fmt.Println("Note: This operation overwrites any existing .ropci.yaml configuration file.\n")

	fmt.Print("Azure Tenant Name or ID (e.g. contoso.onmicrosoft.com): ")
	reader := bufio.NewReader(os.Stdin)

	tenant, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print("*** Error reading tenant information.", err)
		os.Exit(1)
	}

	fmt.Print("Username (e.g. bob@example.org): ")
	user, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print("*** Error reading username.", err)
		os.Exit(1)
	}

	fmt.Println("Nearly done, let's enter the password.")
	fmt.Println("You can leave the password blank if you don't want it stored, and specify the -P flag each time.")
	fmt.Print("Password: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Print("*** Error reading password.", err)
		os.Exit(1)
	}

	// Create the config file
	cfgFile, err := os.Create(".ropci.yaml")
	if err != nil {
		fmt.Printf("*** error creating config file .ropci.yaml: %v\n", err)
	}
	defer cfgFile.Close()

	for _, line := range utils.ReadFileAsStringArray("./templates/config.template") {

		//quick and easy template replacement
		line = strings.Replace(line, "%%USERNAME%%", user, -1)
		line = strings.Replace(line, "%%TENANT%%", tenant, -1)
		line = strings.Replace(line, "%%PASSWORD%%", string(bytes), -1)

		cfgFile.WriteString(line + "\n")
	}

	fmt.Println("\nConfiguration complete.")
	os.Exit(1)
}

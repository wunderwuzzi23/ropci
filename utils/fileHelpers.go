package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"ropci/models"
)

func ReadFileAsStringArray(filename string) []string {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0640)
	if err != nil {
		fmt.Printf("Error: File %s. %s\n", filename, err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func WriteStringToFile(resultsfile string, content string) {
	if resultsfile != "" {
		outfile, err := os.OpenFile(resultsfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer outfile.Close()
		outfile.WriteString(content)
	}
}

func PersistTokenToFile(token *models.Token, filename string) error {

	//Create .token file
	outfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer outfile.Close()

	line, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		fmt.Println("*** Error while marshaling", err)
		return err
	}

	_, err = outfile.WriteString(string(line))
	if err != nil {
		fmt.Printf("*** Error writing token: %s", err)
		return err
	}

	fmt.Printf("Succeeded. Token written to %s.\n", filename)
	return nil
}

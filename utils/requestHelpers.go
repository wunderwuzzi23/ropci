/*
Copyright © 2022 wunderwuzzi23
*/

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	UserAgent string
	Verbose   bool
)

func DoRequest(client *http.Client, graphuri string, apiVersion string, path string, outputFormat string,
	outputFilename string, selectFields []string, search string, retrieveAllRecords bool, reqBody string) {

	//in json mode the entire info is always dumped
	if outputFormat == "json" {
		selectFields = nil
	}

	var sel = ""
	if len(selectFields) > 0 {
		sel := "$select="
		for _, v := range selectFields {
			sel += v + ","
		}
	}

	fullUri := graphuri
	if path != "" {
		fullUri = fullUri + "/" + path
	}

	//some APIs (like Auzre RM, need an API version specified)
	if apiVersion != "" {
		fullUri += "?api-version=" + apiVersion
	} else {
		fullUri = fullUri + "?"
	}

	if sel != "" {
		fullUri = fullUri + sel
	}

	if search != "" {
		fullUri = fullUri + "&$search=\"" + search + "\""
	}

	if path == "search/query" { // search doesn't work well with select
		fullUri = graphuri + "/" + path
	}

	if Verbose {
		fmt.Printf("Request:  %s\n", fullUri)
		if reqBody != "" {
			fmt.Printf("POST Body: %s\n", reqBody)
		}
	}

	getResults(
		client,
		fullUri,
		outputFormat,
		outputFilename,
		selectFields,
		search, //search
		retrieveAllRecords,
		reqBody)

}

func getResults(client *http.Client, fullUri string, outputFormat string, resultsfile string, selectedFields []string, search string, retrieveAllRecords bool, requestBody string) {

	var err error
	var jsonOutfile *os.File

	if outputFormat == "json" && resultsfile != "" {
		jsonOutfile, err = os.OpenFile(resultsfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer jsonOutfile.Close()
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	if outputFormat == "table" {
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(false)
		table.SetHeader(selectedFields)
	}

	count := 0
	var moreData interface{}

	// loop as long as we have new data/uris
	for fullUri != "" {

		var resp *http.Response
		var req *http.Request
		if requestBody == "" { // quick addition to support search/query
			req, _ = http.NewRequest("GET", fullUri, nil)
		} else {
			req, _ = http.NewRequest("POST", fullUri, bytes.NewBufferString(requestBody))
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("User-Agent", UserAgent)
		//fmt.Println("Adding header Prefer")
		//req.Header.Set("Prefer", "outlook.body-content-type=\"text\"") //for mail requests, prefer text over html

		// Request with $search query parameter only works through MSGraph with
		// a special request header: 'ConsistencyLevel: eventual'
		//if search != "" {
		req.Header.Add("ConsistencyLevel", "eventual")
		//}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("*** Error issuing request: ", err.Error())
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("*** Error occured while reading content: " + err.Error())
		}

		if resp.StatusCode < 200 || resp.StatusCode > 204 {
			fmt.Println(string(body))
			return
		}

		if resp.ContentLength == 0 {
			fmt.Println("\nSorry, no results.")
			return
		}

		if Verbose {
			fmt.Println(string(body))
		}

		//parse the returned OData object
		var odata map[string]interface{}
		err = json.Unmarshal(body, &odata)
		if err != nil {
			fmt.Printf("*** Error occured odata: " + err.Error())
		}

		if outputFormat == "table" || outputFormat == "csv" {

			value := odata["value"]

			// Hack to make '/search/query' work - it has a different schema
			// and I don't have time to define models for each schema right now
			if requestBody != "" && strings.Contains(fullUri, "/search/query") {
				value = value.([]interface{})[0].(map[string]interface{})["hitsContainers"].([]interface{})[0].(map[string]interface{})["hits"]
			}

			// if we can't find/cast anything useful, then print the entire body as best effort
			if value == nil {
				fmt.Println(string(body))
				return
			}

			for _, root := range value.([]interface{}) {

				row := []string{}

				switch el := root.(type) {
				case map[string]interface{}:
					for _, v := range selectedFields {
						switch field := el[v].(type) {
						case string:
							row = append(row, field)
						case bool:
							row = append(row, strconv.FormatBool(field))
						case nil:
							row = append(row, "")
						case float64:
							row = append(row, fmt.Sprint(field))
						}
					}

					table.Append(row)
					if Verbose {
						fmt.Println(row)
					}

					count++

				case []interface{}:
					fmt.Printf("List of %d items\n", len(el))
				case string:
					fmt.Printf("%s\n", root)
				default:
					fmt.Printf("*** Encountered unexpected type (%T)", root)
				}
			}
		}

		if outputFormat == "json" {

			count++
			//write to file or stdout
			if resultsfile == "" {
				fmt.Println(string(body))
			} else {
				_, err = jsonOutfile.Write(body)
				if err != nil {
					fmt.Println("*** Error writing result page to file", err)
				}
			}
		}

		moreData = odata["@odata.nextLink"]
		if moreData == nil {
			break
		}

		fullUri = odata["@odata.nextLink"].(string)

		if !retrieveAllRecords {
			break
		}
	}

	// only print the number of pages when writing to file
	// when printing to stdout this allows the user to send the output
	// directly to jq without a parsing error
	if outputFormat == "json" && resultsfile != "" {
		fmt.Printf("Number of pages: %d\n", count)
	}

	if outputFormat == "table" {
		table.Render()
		fmt.Println(tableString.String())
		WriteStringToFile(resultsfile, tableString.String())

		fmt.Printf("Number of items: %d\n", count)
	}

	if moreData != nil {
		fmt.Fprintf(os.Stderr, "*** Only showing limited results, use --all to list everything.\n")
	}

	if resultsfile != "" {
		fmt.Printf("Results written to %s\n", resultsfile)
	}
}

func GenericRequest(client *http.Client, verb string, uri string, requestBody []byte, contentType string) ([]byte, error) {

	//TODO: support adding custom headers here in case it's useful in future
	return requestHelper(client, verb, uri, "", requestBody, contentType)
}

func requestHelper(client *http.Client, verb string, graphuri string, path string, requestBody []byte, contentType string) ([]byte, error) {

	fullUri := graphuri

	if path != "" {
		fullUri = fullUri + "/" + path
	}

	if Verbose {
		fmt.Printf("Request:  %s\n", fullUri)
		fmt.Printf("Verb: %s\n", verb)

		if len(requestBody) > 0 {
			fmt.Printf("POST Body: %v\n", string(requestBody))
		}
	}

	req, err := http.NewRequest(verb, fullUri, bytes.NewBufferString(string(requestBody)))
	if err != nil {
		fmt.Printf("*** Error occured while creating request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("ConsistencyLevel", "eventual")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("*** Error occured while performing request: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 204 {
		fmt.Println(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("*** Error occured while reading content: " + err.Error())
		return nil, err
	}

	if Verbose {
		fmt.Println(resp.Status)
		fmt.Println(string(body))
	}

	return body, nil

}

func Delete(client *http.Client, graphuri string, path string) ([]byte, error) {
	return requestHelper(client, "DELETE", graphuri, path, nil, "")
}
func Post(client *http.Client, graphuri string, path string, requestBody []byte, contentType string) ([]byte, error) {
	return requestHelper(client, "POST", graphuri, path, requestBody, contentType)
}
func Patch(client *http.Client, graphuri string, path string, requestBody []byte, contentType string) ([]byte, error) {
	return requestHelper(client, "PATCH", graphuri, path, requestBody, contentType)
}
func Put(client *http.Client, graphuri string, path string, requestBody []byte, contentType string) ([]byte, error) {
	return requestHelper(client, "PUT", graphuri, path, requestBody, contentType)
}
func Get(client *http.Client, graphuri string, path string, contentType string) ([]byte, error) {
	return requestHelper(client, "GET", graphuri, path, nil, contentType)
}

func GetPrettyJSON(jsonIn string) (string, error) {
	var jsonOut bytes.Buffer
	err := json.Indent(&jsonOut, []byte(jsonIn), "", "  ")
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return jsonOut.String(), nil
}

/*
Copyright © 2022 WUNDERWUZZI23
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"ropci/utils"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var driveUploadLocalFilepath string
var driveContentType string

// driveUploadCmd represents the driveUpload command
// TODO: To support large files (larger then 4MB) implement and upload session
// https://docs.microsoft.com/en-us/graph/api/driveitem-createuploadsession?view=graph-rest-1.0
var driveUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to SharePoint or Drive",
	Long:  `The maximum file size to upload via this API is 4MB. TODO: implement UploadSession for larger files`,
	Run: func(cmd *cobra.Command, args []string) {
		processViper()
		upload()
	},
}

func init() {
	driveCmd.AddCommand(driveUploadCmd)

	driveUploadCmd.Flags().StringVarP(&drivePath, "path", "p", "/", "path of items /")
	driveUploadCmd.Flags().StringVarP(&driveUploadLocalFilepath, "file", "f", "", "local file to upload to drive")
	driveUploadCmd.Flags().StringVarP(&driveContentType, "type", "t", "text", "text or binary")
	driveUploadCmd.MarkFlagRequired("path")
	driveUploadCmd.MarkFlagRequired("file")
}

func upload() {

	if drivePath == "/" {
		fmt.Printf("Please provide a full path and filename, e.g. --path /upload.txt\n")
		return
	}

	fileContents, err := os.ReadFile(driveUploadLocalFilepath)
	if err != nil {
		fmt.Println("*** local file not found", err)
		return
	}

	fullUri := rootGraphUri + "/me/drive/root:" + drivePath + ":/content"

	httpClient := &http.Client{
		Timeout:   2 * time.Second,
		Transport: &utils.UserAgentRoundTripper{},
	}
	_ = context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	fmt.Println("Issuing request: " + fullUri)

	req, err := http.NewRequest("PUT", fullUri, bytes.NewBuffer(fileContents))
	if err != nil {
		fmt.Println("*** Error occured while creating request: " + err.Error())
	}

	var contentType string
	if driveContentType == "binary" {
		contentType = "binary/octet-stream"
	} else {
		contentType = "text/plain"
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", utils.UserAgent)

	resp, err := mainClient.Do(req)
	if err != nil {
		fmt.Println("*** Error occured while performing request: " + err.Error())
	}
	defer resp.Body.Close()
	fmt.Printf("Status %s\n", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("*** Error occured while reading content: " + err.Error())
	}

	fmt.Println(string(body))
}

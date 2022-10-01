package utils

import "fmt"

func GetAzureTokenEndpoint(tenant string) string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenant)
}

func GetAzureDeviceCodeEndpoint(tenant string) string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode", tenant)
}

func GetAzureAuthEndpoint(tenant string) string {
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/auth", tenant)
}

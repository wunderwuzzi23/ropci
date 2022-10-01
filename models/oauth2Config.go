package models

type OAuth2Config struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	Scopes       []string
	//GrantType    string

	TokenEndpoint string
}

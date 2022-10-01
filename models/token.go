package models

import "time"

type Token struct {
	Time         string    `json:"time"`
	DisplayName  string    `json:"display_name"`
	ClientID     string    `json:"client_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IDToken      string    `json:"id_token"`
	Scope        string    `json:"scope"`
	ExpiresIn    int       `json:"expires_in"`
	Expiry       time.Time `json:"expiry"`
	Foci         string    `json:"foci"`
	Result       string    `json:"result"`
	Error        string    `json:"error"`
}

package utils

import "net/http"

// UserAgentRoundTripper helps to ingest a custom User-Agent to the oauth2 library
type UserAgentRoundTripper struct{}

func (*UserAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	req.Header.Set("User-Agent", UserAgent)

	// call the default rounttrip
	response, err := http.DefaultTransport.RoundTrip(req)

	return response, err
}

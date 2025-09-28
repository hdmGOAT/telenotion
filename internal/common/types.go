package common

import "net/http"

type Client struct {
	Token  string
	HTTP *http.Client
}

func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{HTTP: httpClient, Token: token}
}

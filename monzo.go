package monzo

import (
	"io"
	"log"
	"net/http"
)

const ApiUrl = "https://api.monzo.com"

type Client struct {
	http        http.Client
	accessToken string
}

func New(accessToken string) *Client {
	return &Client{
		http:        http.Client{},
		accessToken: accessToken,
	}
}

func (c *Client) Do(method string, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(method, endpoint, body)

	if err != nil {
		log.Print(err)

		return nil, err
	}

	res, err := c.http.Do(req)

	if err != nil {
		log.Print(err)
	}

	return res, err
}

func (c *Client) NewRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, ApiUrl+endpoint, body)

	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	return req, err
}

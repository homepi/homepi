package client

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL string) (*Client, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{
		Timeout:   20 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return &Client{
		BaseURL:   url,
		UserAgent: "HomePi/v0.0.1-cli",
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				Dial:              dialer.Dial,
				ForceAttemptHTTP2: true,
			},
		},
	}, nil
}

func (c *Client) SetAuthToken(token string) *Client {
	c.Token = token
	return c
}

func (c *Client) MakeRequest(data interface{}, endpoint string, method string) error {
	req, err := http.NewRequest(method, c.GetEndpoint(endpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("ApiToken %s", c.Token))
	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	return json.NewDecoder(response.Body).Decode(&data)
}

func (c *Client) GetEndpoint(path string) string {
	return fmt.Sprintf("%s%s", c.BaseURL.String(), path)
}

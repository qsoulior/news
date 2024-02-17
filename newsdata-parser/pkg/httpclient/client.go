package httpclient

import (
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
)

type Client struct {
	*http.Client
	baseURL string
	headers map[string]string
}

func NewClient(opts ...Option) *Client {
	client := &Client{
		Client:  &http.Client{},
		baseURL: "",
		headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Send(method string, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	resultURL, err := urllib.JoinPath(c.baseURL, url)
	if err != nil {
		return nil, fmt.Errorf("url.JoinPath: %w", err)
	}

	req, err := http.NewRequest(method, resultURL, body)
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Do: %w", err)
	}

	return resp, nil
}

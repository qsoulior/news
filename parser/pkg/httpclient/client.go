package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

func New(opts ...Option) *Client {
	client := &Client{
		client:  &http.Client{},
		baseURL: "",
		headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Send(ctx context.Context, method string, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	resultURL := c.baseURL + url
	req, err := http.NewRequestWithContext(ctx, method, resultURL, body)
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Do: %w", err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return c.Send(ctx, http.MethodGet, url, nil, headers)
}

func (c *Client) Head(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return c.Send(ctx, http.MethodHead, url, nil, headers)
}

func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return c.Send(ctx, http.MethodDelete, url, nil, headers)
}

func (c *Client) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return c.Send(ctx, http.MethodPost, url, body, headers)
}

func (c *Client) Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return c.Send(ctx, http.MethodPut, url, body, headers)
}

func (c *Client) Patch(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return c.Send(ctx, http.MethodPatch, url, body, headers)
}

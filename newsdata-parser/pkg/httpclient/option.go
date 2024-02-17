package httpclient

import "time"

type Option func(*Client)

func URL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func Headers(headers map[string]string) Option {
	return func(c *Client) {
		for key, value := range headers {
			c.headers[key] = value
		}
	}
}

func Timeout(d time.Duration) Option {
	return func(c *Client) {
		c.Client.Timeout = d
	}
}

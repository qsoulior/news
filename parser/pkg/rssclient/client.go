package rssclient

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/qsoulior/news/parser/pkg/httpclient"
)

type Client[T any] struct {
	httpclient *httpclient.Client
}

func New[T any](httpclient *httpclient.Client) *Client[T] {
	client := &Client[T]{
		httpclient: httpclient,
	}

	return client
}

func (c *Client[T]) Get(ctx context.Context, name string, url string, headers map[string]string) ([]T, error) {
	resp, err := c.httpclient.Get(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("c.httpclient.Get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp.StatusCode)
	}

	defer resp.Body.Close()

	d := xml.NewDecoder(resp.Body)
	items := make([]T, 0)

	for t, _ := d.Token(); t != nil; t, _ = d.Token() {
		if ctx.Err() != nil {
			return nil, err
		}

		switch tt := t.(type) {
		case xml.StartElement:
			if tt.Name.Local == name {
				var item T
				if err := d.DecodeElement(&item, &tt); err != nil {
					return nil, fmt.Errorf("d.DecodeElement: %w", err)
				}

				items = append(items, item)
			}
		}
	}

	return items, nil
}

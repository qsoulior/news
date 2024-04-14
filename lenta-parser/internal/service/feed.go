package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/DataHenHQ/useragent"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/parser/pkg/httpclient"
	"github.com/qsoulior/news/parser/pkg/rssclient"
	"github.com/rs/zerolog"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	PubDate string   `xml:"pubDate"`
}

type newsFeed struct {
	*news
	url       string
	rssclient *rssclient.Client[Item]
	urlCache  map[string]time.Time
}

func NewNewsFeed(appID string, url string, client *httpclient.Client, logger *zerolog.Logger) *newsFeed {
	log := logger.With().Str("service", "feed").Logger()

	news := &news{
		appID:  appID,
		client: client,
		logger: &log,
	}

	feed := &newsFeed{
		news:      news,
		url:       url,
		rssclient: rssclient.New[Item](client),
		urlCache:  make(map[string]time.Time),
	}

	return feed
}

func (n *newsFeed) Parse(ctx context.Context, query string, page string) ([]entity.News, string, error) {
	urls, err := n.parseURLs(ctx)
	if err != nil {
		return nil, "", err
	}

	news, err := n.parseMany(ctx, urls)
	if err != nil {
		return nil, "", fmt.Errorf("n.parseMany: %w", err)
	}

	// set of output urls
	urlSet := make(map[string]struct{}, len(news))
	for _, item := range news {
		urlSet[item.Link] = struct{}{}
	}

	// delete url that is not in output
	for _, url := range urls {
		if _, ok := urlSet[url.URL]; !ok {
			delete(n.urlCache, url.URL)
		}
	}

	return news, "", nil
}

func (n *newsFeed) parseURLs(ctx context.Context) ([]*newsURL, error) {
	u := n.url + "/rss/news"

	ua, err := useragent.Desktop()
	if err != nil {
		return nil, fmt.Errorf("useragent.Desktop: %w", err)
	}

	items, err := n.rssclient.Get(ctx, "item", u, map[string]string{
		"User-Agent": ua,
	})
	if err != nil {
		return nil, fmt.Errorf("n.rssclient.Get: %w", err)
	}

	// set of current rss urls
	urlSet := make(map[string]struct{}, len(items))

	urls := make([]*newsURL, 0, len(items))
	for _, item := range items {
		url := item.Link
		urlSet[url] = struct{}{}

		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}

		if pd, ok := n.urlCache[url]; !ok || pubDate.After(pd) {
			urls = append(urls, &newsURL{
				URL:         item.Link,
				PublishedAt: pubDate,
			})
			n.urlCache[url] = pubDate
		}

	}

	// delete urls that are not in current rss
	for url := range n.urlCache {
		if _, ok := urlSet[url]; !ok {
			delete(n.urlCache, url)
		}
	}

	return urls, nil
}

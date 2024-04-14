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

const TYPE_ARTICLE = "article"

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	PubDate string   `xml:"pubDate"`
	Type    string   `xml:"type"`
}

type newsFeed struct {
	*news
	url       string
	rssclient *rssclient.Client[Item]
	urlCache  map[string]time.Time
}

func NewNewsFeed(appID string, client *httpclient.Client, url string, logger *zerolog.Logger) *newsFeed {
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

	return news, "", nil
}

func (n *newsFeed) parseURLs(ctx context.Context) ([]string, error) {
	u := n.url + "/export/rss2/archive/index.xml"

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

	// set of rss links
	links := make(map[string]struct{}, len(items))

	urls := make([]string, 0, len(items))
	for _, item := range items {
		if item.Type != TYPE_ARTICLE {
			continue
		}

		link := item.Link
		links[link] = struct{}{}

		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}

		if pd, ok := n.urlCache[link]; !ok || pubDate.After(pd) {
			urls = append(urls, item.Link)
		}

		n.urlCache[link] = pubDate
	}

	// clear cache
	for link := range n.urlCache {
		if _, ok := links[link]; !ok {
			delete(n.urlCache, link)
		}
	}

	return urls, nil
}

package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/internal/service"
	"github.com/rs/zerolog"
)

type NewsConfig struct {
	Logger  *zerolog.Logger
	Service service.News
}

type news struct {
	NewsConfig
}

func NewNews(cfg NewsConfig) *news {
	return &news{cfg}
}

type GetResponse struct {
	Results    []entity.NewsHead `json:"results"`
	Count      int               `json:"count"`
	TotalCount int               `json:"total_count"`
}

func (n *news) getInt(values url.Values, key string) (int, bool) {
	value := values.Get(key)
	if value == "" {
		return 0, false
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}

	return valueInt, true
}

func (n *news) Get(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	query := service.Query{
		Text: values.Get("text"),
	}

	sources := values["sources[]"]
	if len(sources) > 0 {
		query.Sources = make([]string, len(sources))
		copy(query.Sources, sources)
	}

	var opts service.Options
	if skip, ok := n.getInt(values, "skip"); ok {
		opts.SetSkip(skip)
	}

	if limit, ok := n.getInt(values, "limit"); ok {
		opts.SetLimit(limit)
	}

	if sort, ok := n.getInt(values, "sort"); ok {
		opts.SetSort(sort)
	}

	news, count, err := n.Service.GetHead(r.Context(), query, opts)
	if err != nil {
		ErrorJSON(w, "unexpected error while receiving data", http.StatusInternalServerError)
		n.Logger.Error().Err(err).Send()
		return
	}

	if len(news) < 5 && query.Text != "" {
		err := n.Service.Parse(r.Context(), query.Text)
		if err != nil {
			n.Logger.Error().Err(err).Send()
		}
	}

	respData := &GetResponse{
		Results:    news,
		Count:      len(news),
		TotalCount: count,
	}

	EncodeJSON(w, respData, http.StatusOK)
}

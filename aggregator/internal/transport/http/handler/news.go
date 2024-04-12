package handler

import (
	"net/http"
	"strconv"

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

func (n *news) Get(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	query := service.Query{
		Text:   values.Get("text"),
		Source: values.Get("source"),
	}

	var opts service.Options
	skip := values.Get("skip")
	if skip != "" {
		Skip, err := strconv.Atoi(skip)
		if err == nil {
			opts.Skip = Skip
		}
	}

	limit := values.Get("limit")
	if limit != "" {
		Limit, err := strconv.Atoi(limit)
		if err == nil {
			opts.Limit = Limit
		}
	}

	news, count, err := n.Service.GetByQuery(r.Context(), query, opts)
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

	EncodeJSON(w, map[string]any{
		"results":     news,
		"count":       len(news),
		"total_count": count,
	}, http.StatusOK)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/qsoulior/news/aggregator/service"
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

func (n *news) get(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	query := service.Query{
		Title:  values.Get("title"),
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
			opts.Skip = Limit
		}
	}

	news, err := n.Service.GetByQuery(query, opts)
	if err != nil {
		ErrorJSON(w, "unexpected error while receiving data", http.StatusInternalServerError)
		n.Logger.Error().Err(err).Msg("")
		return
	}

	if len(news) < 5 && query.Title != "" {
		err := n.Service.Parse(query.Title)
		if err != nil {
			n.Logger.Error().Err(err).Msg("")
		}
	}

	EncodeJSON(w, news, http.StatusOK)
}

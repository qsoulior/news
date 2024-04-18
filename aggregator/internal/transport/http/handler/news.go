package handler

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/qsoulior/news/aggregator/entity"
	"github.com/qsoulior/news/aggregator/internal/service"
	"github.com/rs/zerolog"
)

type news struct {
	service service.News
}

func NewNews(service service.News) *news {
	return &news{service}
}

type GetResponse struct {
	Results    []entity.NewsHead `json:"results"`
	Skip       uint              `json:"skip"`
	Limit      uint              `json:"limit"`
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

func (n *news) List(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	values := r.URL.Query()
	query := service.Query{
		Text: values.Get("text"),
	}

	sources := values["sources[]"]
	if len(sources) > 0 {
		query.Sources = make([]string, len(sources))
		copy(query.Sources, sources)
	}

	tags := values["tags[]"]
	if len(tags) > 0 {
		query.Tags = make([]string, len(tags))
		copy(query.Tags, tags)
	}

	dateFrom := values.Get("date_from")
	if dateFrom != "" {
		dateFromObj, err := time.Parse(time.DateOnly, dateFrom)
		if err == nil {
			query.DateFrom = &dateFromObj
		}
	}

	dateTo := values.Get("date_to")
	if dateTo != "" {
		dateToObj, err := time.Parse(time.DateOnly, dateTo)
		if err == nil {
			query.DateTo = &dateToObj
		}
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

	news, count, err := n.service.GetHead(r.Context(), query, opts)
	if err != nil {
		ErrorJSON(w, "unexpected error while receiving data", http.StatusInternalServerError)
		logger.Error().Err(err).Send()
		return
	}

	if len(news) < 5 && query.Text != "" {
		err := n.service.SendToParse(r.Context(), query.Text)
		if err != nil {
			logger.Error().Err(err).Send()
		}
	}

	respData := &GetResponse{
		Results:    news,
		Skip:       opts.GetSkip(),
		Limit:      opts.GetLimit(),
		Count:      len(news),
		TotalCount: count,
	}

	EncodeJSON(w, respData, http.StatusOK)
}

func (n *news) Get(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	id := chi.URLParam(r, "id")

	news, err := n.service.Get(r.Context(), id)
	if err != nil {
		ErrorJSON(w, "unexpected error while receiving data", http.StatusInternalServerError)
		logger.Error().Err(err).Send()
		return
	}

	if news == nil {
		ErrorJSON(w, "news with given ID not found", http.StatusNotFound)
		return
	}

	EncodeJSON(w, news, http.StatusOK)
}

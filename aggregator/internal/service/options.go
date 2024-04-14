package service

import "github.com/qsoulior/news/aggregator/internal/repo"

const (
	MaxLimit     = 50
	DefaultLimit = 20
)

type Options struct {
	raw repo.Options
}

func (o *Options) SetSkip(skip int) {
	if skip < 0 {
		o.raw.Skip = 0
	}

	o.raw.Skip = uint(skip)
}

func (o *Options) SetLimit(limit int) {
	if limit <= 0 {
		o.raw.Limit = DefaultLimit
	}

	if limit > MaxLimit {
		o.raw.Limit = MaxLimit
	}

	o.raw.Limit = uint(limit)
}

func (o *Options) SetSort(sort int) {
	if !o.raw.Sort.IsValid() {
		o.raw.Sort = repo.SortDefault
	}

	o.raw.Sort = repo.SortOption(sort)
}

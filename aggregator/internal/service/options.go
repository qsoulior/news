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
		return
	}

	o.raw.Skip = uint(skip)
}

func (o *Options) GetSkip() uint {
	return o.raw.Skip
}

func (o *Options) SetLimit(limit int) {
	if limit <= 0 {
		o.raw.Limit = DefaultLimit
		return
	}

	if limit > MaxLimit {
		o.raw.Limit = MaxLimit
		return
	}

	o.raw.Limit = uint(limit)
}

func (o *Options) GetLimit() uint {
	return o.raw.Limit
}

func (o *Options) SetSort(sort int) {
	if !o.raw.Sort.IsValid() {
		o.raw.Sort = repo.SortDefault
		return
	}

	o.raw.Sort = repo.SortOption(sort)
}

func (o *Options) GetSort() repo.SortOption {
	return o.raw.Sort
}

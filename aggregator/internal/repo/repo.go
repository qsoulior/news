package repo

import (
	"context"

	"github.com/qsoulior/news/aggregator/entity"
)

type News interface {
	Create(ctx context.Context, news entity.News) error
	ReplaceOrCreate(ctx context.Context, news entity.News) error
	CreateMany(ctx context.Context, news []entity.News) error
	GetByID(ctx context.Context, id string) (*entity.News, error)
	GetByQuery(ctx context.Context, query Query, opts Options) ([]entity.NewsHead, int, error)
}

type Query struct {
	Text       string
	Title      bool
	Sources    []string
	Authors    []string
	Tags       []string
	Categories []string
}

type Options struct {
	Limit uint
	Skip  uint
	Sort  SortOption
}

type SortOption uint

const (
	SortPublishedAtDesc SortOption = iota
	SortPublishedAtAsc
	SortRelevanceDesc
	SortRelevanceAsc
	SortDefault = SortPublishedAtDesc
)

func (s SortOption) IsValid() bool {
	return s >= SortPublishedAtDesc && s <= SortRelevanceAsc
}

func (s SortOption) IsPublishedAt() bool {
	return s == SortPublishedAtDesc || s == SortPublishedAtAsc
}

func (s SortOption) IsRelevance() bool {
	return s == SortRelevanceDesc || s == SortRelevanceAsc
}

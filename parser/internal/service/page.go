package service

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/parser/internal/repo"
)

type page struct {
	PageConfig
}

type PageConfig struct {
	Repo repo.Page
}

func NewPage(cfg PageConfig) *page {
	return &page{
		PageConfig: cfg,
	}
}

func (p *page) Get(ctx context.Context) (string, error) {
	page, err := p.Repo.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("p.Repo.Page.Get: %w", err)
	}

	return page, nil

}

func (p *page) Set(ctx context.Context, page string) error {
	err := p.Repo.Update(ctx, page)
	if err != nil {
		return fmt.Errorf("p.Repo.Page.Update: %w", err)
	}

	return nil
}

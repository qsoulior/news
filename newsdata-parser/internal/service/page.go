package service

import (
	"context"
	"fmt"

	"github.com/qsoulior/news/parser/repo"
)

type page struct {
	PageConfig
}

type PageConfig struct {
	repo repo.Page
}

func NewPage(cfg PageConfig) *page {
	return &page{
		PageConfig: cfg,
	}
}

func (p *page) Get() (string, error) {
	page, err := p.repo.Get(context.Background())
	if err != nil {
		return "", fmt.Errorf("p.Repo.Page.Get: %w", err)
	}

	return page, nil

}

func (p *page) Set(page string) error {
	err := p.repo.Update(context.Background(), page)
	if err != nil {
		return fmt.Errorf("p.Repo.Page.Update: %w", err)
	}

	return nil
}

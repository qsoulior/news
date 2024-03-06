package app

import (
	"github.com/qsoulior/news/newsdata-parser/internal/service"
	"github.com/qsoulior/news/parser/app"
)

func Run() {
	news := service.NewNews(service.NewsConfig{})
	page := service.NewPage(service.PageConfig{})
	app.Run(news, page)
}

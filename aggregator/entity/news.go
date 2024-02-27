package entity

import "time"

type News struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Authors     []string  `json:"authors"`
	Tags        []string  `json:"tags"`
	Categories  []string  `json:"categories"`
	Content     string    `json:"content"`
}

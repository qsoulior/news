package entity

import "time"

type News struct {
	Title       string    `json:"title" bson:"title"`
	Link        string    `json:"link" bson:"link"`
	Source      string    `json:"source" bson:"source"`
	PublishedAt time.Time `json:"published_at" bson:"published_at"`
	Authors     []string  `json:"authors" bson:"authors"`
	Tags        []string  `json:"tags" bson:"tags"`
	Categories  []string  `json:"categories" bson:"categories"`
	Content     string    `json:"content" bson:"content"`
}

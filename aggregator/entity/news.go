package entity

import "time"

type NewsHead struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Source      string    `json:"source" bson:"source"`
	PublishedAt time.Time `json:"published_at" bson:"published_at"`
}

type News struct {
	NewsHead   `bson:"inline"`
	Link       string   `json:"link" bson:"link"`
	Authors    []string `json:"authors" bson:"authors"`
	Tags       []string `json:"tags" bson:"tags"`
	Categories []string `json:"categories" bson:"categories"`
	Content    string   `json:"content" bson:"content"`
}

package model

import "time"

type Post struct {
	ID          int64     `db:"id"          json:"id"`
	Slug        string    `db:"slug"         json:"slug"`
	Title       string    `db:"title"        json:"title"`
	Description string    `db:"description"  json:"description"`
	Tags        string    `db:"tags"         json:"tags"` // comma-separated
	Draft       bool      `db:"draft"        json:"draft"`
	PublishDate time.Time `db:"publish_date" json:"publish_date"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updated_at"`
}

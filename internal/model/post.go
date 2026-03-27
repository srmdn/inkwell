package model

import "time"

type Post struct {
	ID          int64     `db:"id"`
	Slug        string    `db:"slug"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Tags        string    `db:"tags"` // comma-separated
	Draft       bool      `db:"draft"`
	PublishDate time.Time `db:"publish_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

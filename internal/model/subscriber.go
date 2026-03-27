package model

import "time"

type Subscriber struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	SubscribedAt time.Time `json:"subscribed_at"`
}

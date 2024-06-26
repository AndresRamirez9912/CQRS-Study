package events

import "time"

type Message interface {
	Type() string
}

type CreatedFeedMessage struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (m CreatedFeedMessage) Type() string {
	return "created_feed"
}

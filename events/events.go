package events

import "events/models"

type EventStore interface {
	Close()
	PublishCreatedFeed(feed *models.Feed) error
	SubscribeCreatedFeed() (<-chan CreatedFeedMessage, error)
	OnCreatedFeed(function func(CreatedFeedMessage)) error
}

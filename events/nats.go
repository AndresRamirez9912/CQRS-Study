package events

import (
	"bytes"
	"encoding/gob"
	"events/models"

	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	conn            *nats.Conn
	feedCreatedSubs *nats.Subscription
	feedCreatedChan chan CreatedFeedMessage
}

func newNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventStore{conn: conn}, nil
}

func (n *NatsEventStore) Close() {
	if n.conn != nil {
		n.conn.Close()
	}

	if n.feedCreatedSubs != nil {
		n.feedCreatedSubs.Unsubscribe()
	}

	close(n.feedCreatedChan)
}

func (n *NatsEventStore) PublishCreatedFeed(feed *models.Feed) error {
	msg := &CreatedFeedMessage{
		Id:          feed.Id,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	encodedMsg, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}

	return n.conn.Publish(msg.Type(), encodedMsg)
}

func (n *NatsEventStore) OnCreatedFeed(function func(CreatedFeedMessage)) error {
	msg := CreatedFeedMessage{}
	subs, err := n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMessage(m.Data, &msg)
		function(msg)
	})
	n.feedCreatedSubs = subs
	return err
}

func (n *NatsEventStore) SubscribeCreatedFeed() (<-chan CreatedFeedMessage, error) {
	createdMessage := CreatedFeedMessage{}
	n.feedCreatedChan = make(chan CreatedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	subs, err := n.conn.ChanSubscribe(createdMessage.Type(), ch)
	if err != nil {
		return nil, err
	}
	n.feedCreatedSubs = subs
	go func() {
		for {
			select {
			case msg := <-ch:
				n.decodeMessage(msg.Data, &createdMessage)
				n.feedCreatedChan <- createdMessage
			}
		}
	}()
	return (<-chan CreatedFeedMessage)(n.feedCreatedChan), nil
}

func (n *NatsEventStore) encodeMessage(message Message) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := gob.NewEncoder(&buffer).Encode(message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func (n *NatsEventStore) decodeMessage(data []byte, inter interface{}) error {
	bytes := bytes.NewBuffer(data)
	return gob.NewDecoder(bytes).Decode(inter)
}

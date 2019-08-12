package events

import (
	"context"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"
)

const (
	ORMChangeEvent = "com.graphql.orm.change"
)

type EventController struct {
	client *cloudevents.Client
}

func NewEventController() (ec EventController, err error) {
	URL := os.Getenv("EVENT_TRANSPORT_URL")
	var _client *cloudevents.Client
	if URL != "" {
		t, tErr := cloudevents.NewHTTPTransport(
			cloudevents.WithTarget(URL),
			cloudevents.WithBinaryEncoding(),
		)
		err = tErr
		if err != nil {
			return
		}

		client, cErr := cloudevents.NewClient(t)
		err = cErr
		if err != nil {
			return
		}
		log.Printf("Created cloudevents client with target %s", URL)
		_client = &client
	}
	ec = EventController{_client}
	return
}

func (c *EventController) send(ctx context.Context, e cloudevents.Event) error {
	if c.client == nil {
		return nil
	}
	_, err := (*c.client).Send(ctx, e)
	return err
}

// SendEvent ...
func (c *EventController) SendEvent(ctx context.Context, e *Event) (err error) {
	event := cloudevents.NewEvent()
	event.SetID(e.ID)
	event.SetType(ORMChangeEvent)
	event.SetSource("http://graphql-orm/graphql")
	err = event.SetData(e)
	if err != nil {
		return
	}

	err = c.send(ctx, event)
	return
}

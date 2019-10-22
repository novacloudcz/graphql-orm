package events

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go"
)

const (
	ORMChangeEvent = "com.graphql.orm.change"
)

type EventController struct {
	clients map[string]cloudevents.Client
	debug   bool
}

func NewEventController() (ec EventController, err error) {
	URLs := getENVArray("EVENT_TRANSPORT_URL")
	_clients := map[string]cloudevents.Client{}
	for _, URL := range URLs {
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
			_clients[URL] = client
		}
	}
	debug := os.Getenv("DEBUG") == "true"
	ec = EventController{clients: _clients, debug: debug}
	return
}

func (c *EventController) send(ctx context.Context, e cloudevents.Event) error {
	for URL, client := range c.clients {
		if _, err := client.Send(ctx, e); err != nil {
			if c.debug {
				fmt.Printf("received cloudevents error %s from server %s\n", err.Error(), URL)
			}
			return err
		}
	}
	return nil
}

// SendEvent ...
func (c *EventController) SendEvent(ctx context.Context, e *Event) (err error) {
	if len(c.clients) == 0 {
		return
	}
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

func getENVArray(name string) []string {
	arr := []string{}

	val := os.Getenv(name)
	if val != "" {
		arr = append(arr, strings.Split(val, ",")...)
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%s_%d", name, i)
		sval := os.Getenv(key)
		if sval != "" {
			arr = append(arr, sval)
		}
	}

	return arr
}

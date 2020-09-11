package templates

var EventsController = `package gen

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	cloudeventsaws "github.com/jakubknejzlik/cloudevents-aws-transport"
	"github.com/novacloudcz/graphql-orm/events"
)

const (
	ORMChangeEvent = "com.graphql.orm.change"
)

type EventController struct {
	clients map[string]cloudevents.Client
	debug   bool
	source  string
}

func NewEventController() (ec EventController, err error) {
	source := os.Getenv("EVENT_TRANSPORT_SOURCE")
	if source == "" {
		hostname, _err := os.Hostname()
		if err != nil {
			err = _err
			return
		}
		source = "http://" + hostname + "/graphql"
	}

	URLs := getENVArray("EVENT_TRANSPORT_URL")
	_clients := map[string]cloudevents.Client{}
	for _, URL := range URLs {
		if URL != "" {
			t, tErr := transportForURL(URL)
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
	ec = EventController{clients: _clients, debug: debug, source: source}

	log.Printf("Created EventController with source %s", source)

	return
}

func (c *EventController) send(ctx context.Context, e cloudevents.Event) error {
	for URL, client := range c.clients {
		if _, _, err := client.Send(ctx, e); err != nil {
			if c.debug {
				fmt.Printf("received cloudevents error %s from server %s\n", err.Error(), URL)
			}
			return err
		}
	}
	return nil
}

// SendEvent ...
func (c *EventController) SendEvent(ctx context.Context, e *events.Event) (err error) {
	if len(c.clients) == 0 {
		return
	}
	event := cloudevents.NewEvent()
	event.SetID(e.ID)
	event.SetType(ORMChangeEvent)
	event.SetSource(c.source)
	event.SetTime(e.Date)
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

func transportForURL(URL string) (t transport.Transport, err error) {

	if strings.HasPrefix(URL, "arn:aws:sns") {
		t, err = cloudeventsaws.NewSNSTransport(URL)
		return
	}
	if strings.HasPrefix(URL, "arn:aws:events") {
		t, err = cloudeventsaws.NewEventBridgeTransport(URL)
		return
	}

	u, err := url.Parse(URL)
	if err != nil {
		return
	}
	switch u.Scheme {
	case "http", "https":
		t, err = cloudevents.NewHTTPTransport(
			cloudevents.WithTarget(URL),
			cloudevents.WithBinaryEncoding(),
		)
	case "sqs+https":
		u.Scheme = "https"
		t, err = cloudeventsaws.NewSQSTransport(u.String())
	default:
		err = fmt.Errorf("unknown scheme %s", u.Scheme)

	}
	return
}
`

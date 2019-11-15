package events

import (
	"encoding/json"
	"time"

	uuid "github.com/gofrs/uuid"
)

// EventType ...
type EventType string

const (
	// EventTypeCreated ...
	EventTypeCreated = "CREATED"
	// EventTypeUpdated ...
	EventTypeUpdated = "UPDATED"
	// EventTypeDeleted ...
	EventTypeDeleted = "DELETED"
)

type EventDataValue interface{}

// EventChange ...
type EventChange struct {
	Name     string `json:"name"`
	OldValue string `json:"oldValue"`
	NewValue string `json:"newValue"`
}

func (ec *EventChange) SetOldValue(value interface{}) error {
	data, err := json.Marshal(value)
	if err == nil {
		ec.OldValue = string(data)
	}
	return err
}
func (ec *EventChange) OldValueAs(data interface{}) error {
	return json.Unmarshal([]byte(ec.OldValue), data)
}
func (ec *EventChange) SetNewValue(value interface{}) error {
	data, err := json.Marshal(value)
	if err == nil {
		ec.NewValue = string(data)
	}
	return err
}
func (ec *EventChange) NewValueAs(data interface{}) error {
	return json.Unmarshal([]byte(ec.NewValue), data)
}

type EventMetadata struct {
	Type        EventType `json:"type"`
	Cursor      string    `json:"cursor"`
	Entity      string    `json:"entity"`
	EntityID    string    `json:"entityId"`
	Date        time.Time `json:"date"`
	PrincipalID *string   `json:"principalId"`
}

// Event ...
type Event struct {
	EventMetadata
	ID      string         `json:"id"`
	Changes []*EventChange `json:"changes"`
}

// NewEvent ...
func NewEvent(meta EventMetadata) Event {
	return Event{
		EventMetadata: meta,
		ID:            uuid.Must(uuid.NewV4()).String(),
		Changes:       []*EventChange{},
	}
}

// HasChangedColumn check if given event has changes on specific column
func (e Event) HasChangedColumn(c string) bool {
	for _, col := range e.ChangedColumns() {
		if col == c {
			return true
		}
	}
	return false
}

// ChangedColumns returns list of names of changed columns
func (e Event) ChangedColumns() []string {
	columns := []string{}

	for _, change := range e.Changes {
		columns = append(columns, change.Name)
	}

	return columns
}

func (e *Event) Change(column string) (ec *EventChange) {
	for _, c := range e.Changes {
		if c.Name == column {
			ec = c
			break
		}
	}
	return
}

// AddNewValue ...
func (e *Event) AddNewValue(column string, v EventDataValue) {
	change := e.Change(column)
	if change == nil {
		c := EventChange{Name: column}
		change = &c
		e.Changes = append(e.Changes, change)
	}
	if err := change.SetNewValue(v); err != nil {
		panic("failed to set new value" + err.Error())
	}
}

// AddOldValue ...
func (e *Event) AddOldValue(column string, v EventDataValue) {
	change := e.Change(column)
	if change == nil {
		c := EventChange{Name: column}
		change = &c
		e.Changes = append(e.Changes, change)
	}
	if err := change.SetOldValue(v); err != nil {
		panic("failed to set new value" + err.Error())
	}
}

package events

import (
	"time"

	uuid "github.com/satori/go.uuid"
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
	Name     string          `json:"name"`
	OldValue *EventDataValue `json:"oldValue"`
	NewValue *EventDataValue `json:"newValue"`
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
	ID        string                     `json:"id"`
	Changes   []*EventChange             `json:"changes"`
	OldValues map[string]*EventDataValue `json:"oldValues"`
	NewValues map[string]*EventDataValue `json:"newValues"`
}

// NewEvent ...
func NewEvent(meta EventMetadata) Event {
	return Event{
		EventMetadata: meta,
		ID:            uuid.Must(uuid.NewV4()).String(),
		Changes:       []*EventChange{},
		OldValues:     map[string]*EventDataValue{},
		NewValues:     map[string]*EventDataValue{},
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

// OldValue returns old value for column
func (e Event) OldValue(c string) (*EventDataValue, bool) {
	v, ok := e.OldValues[c]
	return v, ok
}

// NewValue returns new value for column
func (e Event) NewValue(c string) (*EventDataValue, bool) {
	v, ok := e.NewValues[c]
	return v, ok
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
func (e *Event) AddNewValue(column string, value interface{}) {
	v := EventDataValue(value)
	e.NewValues[column] = &v
	change := e.Change(column)
	if change == nil {
		c := EventChange{Name: column}
		change = &c
		e.Changes = append(e.Changes, change)
	}
	change.NewValue = &v
}

// AddOldValue ...
func (e *Event) AddOldValue(column string, value interface{}) {
	v := EventDataValue(value)
	e.NewValues[column] = &v
	change := e.Change(column)
	if change == nil {
		c := EventChange{Name: column}
		change = &c
		e.Changes = append(e.Changes, change)
	}
	change.OldValue = &v
}

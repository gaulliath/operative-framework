package session

import (
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"time"
)

type Events []*Event

type Event struct {
	EventId string
	Type    string
	Value   string
	Date    time.Time
}

func (s *Session) NewEvent(t string, value string) *Event {
	event := &Event{
		EventId: "E_" + ksuid.New().String(),
		Type:    t,
		Value:   value,
		Date:    time.Now(),
	}
	s.Events = append(s.Events, event)
	return event
}

func (s *Session) GetEvent(id string) (*Event, error) {
	for _, event := range s.Events {
		if event.EventId == id {
			return event, nil
		}
	}
	return nil, errors.New("Event '" + id + "' as not found.")
}

func (s *Session) GetEvents() Events {
	return s.Events
}

func (s *Session) DeleteEvent(id string) {
	newEvents := Events{}
	for _, event := range s.Events {
		if event.EventId != id {
			newEvents = append(newEvents, event)
		}
	}
	s.Events = newEvents
}

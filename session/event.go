package session

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"strings"
	"time"
)

type Events []*Event

type Event struct {
	EventId string      `json:"event_id"`
	Type    string      `json:"type"`
	Value   interface{} `json:"raw"`
	JSON    interface{} `json:"json"`
	Date    time.Time   `json:"date"`
}

const (
	EXEC_COMMAND    = "EXEC_COMMAND"
	MONITOR_MATCH   = "MONITOR_MATCH"
	TARGET_ADD      = "TARGET_ADD"
	TARGET_LINK     = "TARGET_LINK"
	TAG_ADD         = "TAG_ADD"
	ERROR_ANALYTICS = "ERROR_ANALYTICS"
	ERROR_GENERIC   = "ERROR_GENERIC"
	RESULTS_ADD     = "RESULTS_ADD"
	RESULTS_SYNC    = "RESULTS_SYNC"
	WEBHOOK_SEND    = "WEBHOOK_SEND"
)

func (s *Session) NewEvent(t string, value interface{}) *Event {

	jsonify, _ := json.Marshal(value)

	event := &Event{
		EventId: "E_" + ksuid.New().String(),
		Type:    t,
		Value:   value,
		JSON:    string(jsonify),
		Date:    time.Now(),
	}
	s.Events = append(s.Events, event)

	for _, wh := range s.WebHooks {
		if wh.GetStatus() {
			for _, e := range wh.Events {
				if strings.ToLower(e) == strings.ToLower(t) {
					err := s.SendToWebHook(wh, event)
					if err != nil {
						s.Stream.Error(err.Error())
						s.NewEvent(ERROR_GENERIC, err.Error())
						return event
					}
					s.NewEvent(WEBHOOK_SEND, wh)
				}
			}
		}
	}

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

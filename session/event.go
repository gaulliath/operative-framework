package session

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"strings"
	"time"
)

type Events struct {
	Watcher []string `json:"notifier"`
	Lists   []*Event `json:"lists"`
}

type Event struct {
	EventId string      `json:"event_id"`
	Type    string      `json:"type"`
	Value   interface{} `json:"raw"`
	JSON    interface{} `json:"json"`
	Date    time.Time   `json:"date"`
}

const (
	EXEC_COMMAND     = "EXEC_COMMAND"
	MONITOR_MATCH    = "MONITOR_MATCH"
	TARGET_ADD       = "TARGET_ADD"
	TARGET_LINK      = "TARGET_LINK"
	TAG_ADD          = "TAG_ADD"
	ERROR_GENERIC    = "ERROR_GENERIC"
	ERROR_MODULE     = "ERROR_MODULE"
	RESULTS_ADD      = "RESULTS_ADD"
	RESULTS_SYNC     = "RESULTS_SYNC"
	WEBHOOK_SEND     = "WEBHOOK_SEND"
	TRACKER_FOUND    = "TRACKER_FOUND"
	TRACKER_SELECTED = "TRACKER_SELECTED"
	TRACKER_UNKNOWN  = "TRACKER_UNKNOWN"
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

	if s.HasWatcher(t) {
		s.NewNotification("New event '" + strings.ToLower(t) + "' monitored ")
	}
	s.Events.Lists = append(s.Events.Lists, event)

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

// Checking if watcher as been set in session
func (s *Session) HasWatcher(t string) bool {
	for _, watcher := range s.Events.Watcher {
		if strings.ToLower(watcher) == strings.ToLower(t) {
			return true
		}
	}
	return false
}

func (s *Session) GetEvent(id string) (*Event, error) {
	for _, event := range s.Events.Lists {
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
	for _, event := range s.Events.Lists {
		if event.EventId != id {
			newEvents.Lists = append(newEvents.Lists, event)
		}
	}
	s.Events = newEvents
}

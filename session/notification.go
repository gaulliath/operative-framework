package session

import (
	"errors"
	"github.com/segmentio/ksuid"
	"time"
)

type Notification struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	IsRead    bool   `json:"as_read"`
	CreatedAt time.Time
}

// Set notification with read status
func (n *Notification) Read() {
	n.IsRead = true
	return
}

// Create a new notification
func (s *Session) NewNotification(text string) *Notification {

	notification := Notification{
		Id:        "NTF_" + ksuid.New().String(),
		Text:      text,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	s.Notifications = append(s.Notifications, &notification)
	return &notification
}

// Get a notification
func (s *Session) GetNotification(id string) (*Notification, error) {

	for _, notification := range s.Notifications {
		if notification.Id == id {
			return notification, nil
		}
	}

	return nil, errors.New("Notification '" + id + "' not found")
}

// Count unread session notification
func (s *Session) CountUnReadNotifications() int {
	count := 0
	for _, notification := range s.Notifications {
		if !notification.IsRead {
			count = count + 1
		}
	}
	return count
}

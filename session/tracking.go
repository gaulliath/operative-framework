package session

import (
	"github.com/segmentio/ksuid"
	"sort"
	"time"
)

type Tracking struct {
	Id          string             `json:"id"`
	Session     *Session           `json:"-"`
	Position    TrackingPosition   `json:"position"`
	Identifier  string             `json:"identifier"`
	Description string             `json:"description"`
	Memories    []TrackingPosition `json:"memories"`
	IsHistory   bool               `json:"is_history"`
	Time        time.Time          `json:"time"`
}

type TrackingPosition struct {
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
	Time      time.Time `json:"time"`
}

func (s *Session) AddTracker(tracker Tracking) Tracking {
	tracker.Id = ksuid.New().String()
	tracker.Session = s
	tracker.Memories = []TrackingPosition{}
	tracker.Position.Time = time.Now()
	tracker.Time = time.Now()
	return s.AddOrFirstTracker(tracker)
}

func (s *Session) AddOrFirstTracker(t Tracking) Tracking {
	for _, tracker := range s.Tracker {
		if tracker.Identifier == t.Identifier && tracker.IsHistory == false {
			if tracker.Position.Longitude == t.Position.Longitude &&
				tracker.Position.Latitude == t.Position.Latitude {
				return *tracker
			}

			for _, memory := range tracker.Memories {
				t.Memories = append(t.Memories, memory)
			}
			t.Memories = append(t.Memories, tracker.Position)
			tracker.IsHistory = true
		}
	}
	s.Tracker = append(s.Tracker, &t)
	return t
}

func CreateTrackerFromValue(s *Session, lat string, lng string, id string, description string) *Tracking {
	return &Tracking{
		Session: s,
		Position: TrackingPosition{
			Latitude:  lat,
			Longitude: lng,
			Time:      time.Now(),
		},
		Identifier:  id,
		Description: description,
		Time:        time.Now(),
	}
}

func (t *Tracking) Synchronize() {
	trackers := t.Session.Tracker
	for _, tracker := range trackers {
		if tracker.Identifier == t.Identifier {
			if tracker.Position.Longitude == t.Position.Longitude &&
				tracker.Position.Latitude == t.Position.Latitude {
				continue
			}
			for _, memory := range tracker.Memories {
				t.Memories = append(t.Memories, memory)
			}
			t.Memories = append(t.Memories, tracker.Position)
			tracker.IsHistory = true
		}
	}
}

func (t *Tracking) GetLatitude() string {
	return t.Position.Latitude
}

func (t *Tracking) GetLongitude() string {
	return t.Position.Longitude
}

func (t *Tracking) GetIdentifier() string {
	return t.Identifier
}

func (t *Tracking) GetDescription() string {
	return t.Description
}

func (t *Tracking) GetMemories() []TrackingPosition {
	return t.Memories
}

func (t *Tracking) GetTime() time.Time {
	return t.Time
}

func (t *Tracking) GetHistories() []*Tracking {
	trackers := t.Session.Tracker
	sort.Slice(trackers, func(i, j int) bool {
		return trackers[i].Time.After(trackers[j].Time)
	})
	return trackers
}

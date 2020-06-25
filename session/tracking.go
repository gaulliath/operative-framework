package session

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/graniet/go-pretty/table"
	"github.com/segmentio/ksuid"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type Tracking struct {
	Id          string             `json:"id"`
	Session     *Session           `json:"-"`
	Position    TrackingPosition   `json:"position"`
	Identifier  string             `json:"identifier"`
	Description string             `json:"description"`
	Picture     string             `json:"picture"`
	Memories    []TrackingPosition `json:"memories"`
	IsHistory   bool               `json:"is_history"`
	Time        time.Time          `json:"time"`
}

type Bounds struct {
	Distance  float64  `json:"distance"`
	NorthWest Position `json:"north_west"`
	SouthWest Position `json:"south_west"`
	NorthEast Position `json:"north_east"`
	SouthEast Position `json:"south_east"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TrackingPosition struct {
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
	Bounds    Bounds    `json:"bounds"`
	Time      time.Time `json:"time"`
}

func (s *Session) GetTrackingUrlWithParam() string {
	return "http://" + s.Config.Tracker.Host + ":" + s.Config.Tracker.Port + "?api=" + s.Config.Api.Host + "&port=" + s.Config.Api.Port
}

func (s *Session) AddTracker(tracker Tracking) Tracking {
	tracker.Id = ksuid.New().String()
	tracker.Session = s
	tracker.Memories = []TrackingPosition{}
	tracker.Position.Time = time.Now()
	tracker.Time = time.Now()
	return s.AddOrFirstTracker(tracker)
}

func (s *Session) ResetPosition() {
	s.Tracker.Position = []Position{}
	return
}

func (s *Session) SetPosition(position Position) {
	s.Tracker.Position = append(s.Tracker.Position, position)
	return
}

func (s *Session) GetTracker(trackerId string) (*Tracking, error) {
	for _, track := range s.Tracker.Tracked {
		if track.Id == trackerId {
			s.NewEvent(TRACKER_FOUND, track)
			return track, nil
		}
	}

	s.NewEvent(TRACKER_UNKNOWN, trackerId)
	return nil, errors.New("Tracker identifier not found")
}

func (s *Session) AddOrFirstTracker(t Tracking) Tracking {
	for _, tracker := range s.Tracker.Tracked {
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

	t.Position.ToBounds(140)
	s.Tracker.Tracked = append(s.Tracker.Tracked, &t)

	s.Tracker.Selected = s.GetTrackerBestMover()
	s.SetPositionFromTracker(s.Tracker.Selected)

	return t
}

func (s *Session) GetLastTracker() (*Tracking, error) {
	if len(s.Tracker.Tracked) > 0 {
		return s.Tracker.Tracked[len(s.Tracker.Tracked)-1], nil
	}
	return nil, errors.New("No tracker found.")
}

func (s *Session) SetPositionFromTracker(tracker *Tracking) {
	s.ResetPosition()
	s.Tracker.Position = append(s.Tracker.Position, tracker.Position.Bounds.SouthEast)
	s.Tracker.Position = append(s.Tracker.Position, tracker.Position.Bounds.SouthWest)
	s.Tracker.Position = append(s.Tracker.Position, tracker.Position.Bounds.NorthEast)
	s.Tracker.Position = append(s.Tracker.Position, tracker.Position.Bounds.NorthWest)
	return
}

func (s *Session) GetTrackerBestMover() *Tracking {
	trackers := s.Tracker.Tracked
	sort.Slice(trackers, func(i, j int) bool {
		newTime := trackers[i].Time.Add((10 * time.Minute))
		if time.Now().After(newTime) {
			return false
		}

		return len(trackers[i].Memories) > len(trackers[j].Memories)
	})
	if len(trackers) < 1 {
		return &Tracking{}
	}
	return trackers[0]
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
	trackers := t.Session.Tracker.Tracked
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
	trackers := t.Session.Tracker.Tracked
	sort.Slice(trackers, func(i, j int) bool {
		return trackers[i].Time.After(trackers[j].Time)
	})
	return trackers
}

func (tp *TrackingPosition) ToBounds(m float64) {
	lat, _ := strconv.ParseFloat(tp.Latitude, 64)
	lng, _ := strconv.ParseFloat(tp.Longitude, 64)

	var latAccuracy float64
	var lngAccuracy float64

	latAccuracy = (180 * m) / 40075017
	lngAccuracy = float64(latAccuracy) / math.Cos((math.Pi/180)*lat)
	t2lat := lat + float64(latAccuracy)

	tp.Bounds.Distance = m

	tp.Bounds.NorthWest.Latitude = t2lat
	tp.Bounds.NorthWest.Longitude = lng - lngAccuracy

	tp.Bounds.NorthEast.Latitude = t2lat
	tp.Bounds.NorthEast.Longitude = lng + lngAccuracy

	tp.Bounds.SouthWest.Latitude = lat - float64(latAccuracy)
	tp.Bounds.SouthWest.Longitude = lng - lngAccuracy

	tp.Bounds.SouthEast.Latitude = lat - float64(latAccuracy)
	tp.Bounds.SouthEast.Longitude = lng + lngAccuracy
	return
}

func (t *Tracking) ViewPositions() {
	tbl := t.Session.Stream.GenerateTable()
	tbl.SetOutputMirror(os.Stdout)
	tbl.SetAllowedColumnLengths([]int{40, 30, 30})
	tbl.AppendHeader(table.Row{
		"LATITUDE",
		"LONGITUDE",
		"TIME",
	})

	tbl.AppendRow(table.Row{
		t.Position.Latitude,
		t.Position.Longitude,
		t.Position.Time.Format("2006-01-02 15:04:05"),
	})

	for _, pos := range t.Memories {
		tbl.AppendRow(table.Row{
			pos.Latitude,
			pos.Longitude,
			pos.Time.Format("2006-01-02 15:04:05"),
		})
	}
	t.Session.Stream.Render(tbl)
}

func (s *Session) GetTrackerRouter() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./gui/tracker/static/"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./gui/tracker/index.html")
	})
	return router
}

func (s *Session) ServeTrackerGUI() {
	err := s.Tracker.Server.ListenAndServe()
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}
	return
}

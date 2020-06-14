package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/graniet/operative-framework/session"
	"net/http"
	"sort"
)

func (api *ARestFul) Trackers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var trackers []session.Tracking
	for _, tracker := range api.sess.Tracker.Tracked {
		if tracker.IsHistory == true {
			continue
		}
		trackers = append(trackers, *tracker)
	}

	message := api.Core.PrintData("request executed", false, trackers)
	_ = json.NewEncoder(w).Encode(message)
	return
}

func (api *ARestFul) GetPositions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	message := api.Core.PrintData("request executed", false, api.sess.Tracker.Position)
	_ = json.NewEncoder(w).Encode(message)
	return
}

func (api *ARestFul) GetMovers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	trackers := api.sess.Tracker.Tracked
	sort.Slice(trackers, func(i, j int) bool {
		return len(trackers[i].Memories) > len(trackers[j].Memories)
	})

	message := api.Core.PrintData("request executed", false, trackers)
	_ = json.NewEncoder(w).Encode(message)

	return
}

func (api *ARestFul) GetBestMovers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	message := api.Core.PrintData("request executed", false, api.sess.GetTrackerBestMover())
	_ = json.NewEncoder(w).Encode(message)

	return
}

func (api *ARestFul) Tracker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var t session.Tracking
	params := mux.Vars(r)
	identifier := params["identifier"]
	for _, tracker := range api.sess.Tracker.Tracked {
		if tracker.IsHistory == true {
			continue
		}

		if tracker.Id == identifier {
			t = *tracker
		}
	}

	message := api.Core.PrintData("request executed", false, t)
	_ = json.NewEncoder(w).Encode(message)
	return
}

func (api *ARestFul) PutTracker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var tracker session.Tracking

	err := json.NewDecoder(r.Body).Decode(&tracker)
	if err != nil {
		message := api.Core.PrintMessage(err.Error(), true)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	tracker = api.sess.AddTracker(tracker)
	message := api.Core.PrintData("request executed", false, tracker)
	_ = json.NewEncoder(w).Encode(message)
	return
}

package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/graniet/operative-framework/session"
	"net/http"
)

func (api *ARestFul) Trackers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var trackers []session.Tracking
	for _, tracker := range api.sess.Tracker {
		if tracker.IsHistory == true {
			continue
		}
		trackers = append(trackers, *tracker)
	}

	message := api.Core.PrintData("request executed", false, trackers)
	_ = json.NewEncoder(w).Encode(message)
	return
}

func (api *ARestFul) Tracker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t session.Tracking
	params := mux.Vars(r)
	identifier := params["identifier"]
	for _, tracker := range api.sess.Tracker {
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
	var tracker session.Tracking

	err := json.NewDecoder(r.Body).Decode(&tracker)
	if err != nil {
		message := api.Core.PrintMessage("We can't parse tracker values", true)
		_ = json.NewEncoder(w).Encode(message)
		return
	}

	tracker = api.sess.AddTracker(tracker)
	message := api.Core.PrintData("request executed", false, tracker)
	_ = json.NewEncoder(w).Encode(message)
	return
}

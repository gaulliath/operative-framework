package api

import (
	"encoding/json"
	"github.com/graniet/operative-framework/session"
	"net/http"
)

func (api *ARestFul) Intervals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var intervals []session.Interval
	for _, itv := range api.sess.Interval {
		intervals = append(intervals, *itv)
	}
	message := api.Core.PrintData("requests executed", false, intervals)
	_ = json.NewEncoder(w).Encode(message)
}

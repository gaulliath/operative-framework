package session

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const (
	ONLY_SCREEN       = "ONLY_SERVER"
	ONLY_SERVER       = "ONLY_SERVER"
	SCREEN_AND_SERVER = "SCREEN_AND_SERVER"
)

func (s *Session) PushResultsToGate(results []*TargetResults, afterTime time.Time) {
	for _, result := range results {
		if result.CreatedAt.After(afterTime) {
			j, err := json.Marshal(result)
			if err == nil {
				client := http.Client{}
				_, err := client.Post(s.Config.Gate.GateUrl, "application/json", bytes.NewBuffer(j))
				if err != nil {
					s.Stream.Error(err.Error())
				}
			}
		}
	}
}

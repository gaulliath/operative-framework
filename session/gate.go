package session

import (
	"time"
)

func (s *Session) PushResultsToGate(results []*TargetResults, afterTime time.Time) {
	for _, result := range results {
		if result.CreatedAt.After(afterTime) {
			// Push to gate.
		}
	}
}

package session

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Source struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

func (s *Session) FromSourceFile() error {

	var sources []Source
	file, err := os.Open(s.SourceFile)
	if err != nil {
		return errors.New(err.Error())
	}

	source, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New(err.Error())
	}

	_ = json.Unmarshal(source, &sources)

	if len(sources) > 0 {
		for _, source := range sources {
			content, _ := json.Marshal(source.Content)

			switch strings.ToLower(source.Type) {
			case "interval":
				var interval Interval
				_ = json.Unmarshal(content, &interval)
				interval.SetId()
				interval.SetSession(s)
				interval.SetTimeType("minute")
				s.Interval = append(s.Interval, &interval)
				interval.Up()
				break
			case "monitor":
				var monitor Monitor
				_ = json.Unmarshal(content, &monitor)
				monitor.SetId()
				monitor.SetSession(s)
				monitor.CreatedAt = time.Now()
				s.Monitors = append(s.Monitors, &monitor)
				monitor.Up()
				break
			}
		}
	}
	return nil
}

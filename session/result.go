package session

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type TargetResults struct {
	Id         int             `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionId  int             `json:"-" gorm:"session_id"`
	ModuleName string          `json:"module_name"`
	ResultId   string          `json:"result_id" gorm:"primary_key:yes;column:result_id"`
	TargetId   string          `json:"target_id" gorm:"target_id"`
	Header     string          `json:"key" gorm:"result_header"`
	Value      string          `json:"value" gorm:"result_value"`
	ToJSON     json.RawMessage `json:"to_json"`
	Notes      []Note
	Auxiliary  []string  `json:"-" sql:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *Session) GetResultsAfter(results []*TargetResults, afterTime time.Time) []*TargetResults {
	var scopes []*TargetResults
	for _, result := range results {
		if result.CreatedAt.After(afterTime) {
			scopes = append(scopes, result)
		}
	}

	return scopes
}

func (s *Session) GetResultsBefore(results []*TargetResults, beforeTime time.Time) []*TargetResults {
	var scopes []*TargetResults
	for _, result := range results {
		if result.CreatedAt.Before(beforeTime) {
			scopes = append(scopes, result)
		}
	}

	return scopes
}

func (r *TargetResults) Bytes() ([]byte, error) {
	separator := base64.StdEncoding.EncodeToString([]byte(";operativeframework;"))[0:5]

	format := make(map[string]interface{})
	keys := strings.Split(r.Header, separator)
	values := strings.Split(r.Value, separator)

	if len(keys) > 0 {
		for count, key := range keys {
			if ok := values[count]; ok != "" {
				format[key] = values[count]
			} else {
				format[key] = false
			}
		}
	}

	content, err := json.Marshal(format)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (r *TargetResults) JSON() json.RawMessage {
	separator := base64.StdEncoding.EncodeToString([]byte(";operativeframework;"))[0:5]

	format := make(map[string]interface{})
	keys := strings.Split(r.Header, separator)
	values := strings.Split(r.Value, separator)

	if len(keys) > 0 {
		for count, key := range keys {
			if ok := values[count]; ok != "" {
				format[key] = values[count]
			} else {
				format[key] = false
			}
		}
	}

	content, err := json.Marshal(format)
	if err != nil {
		return nil
	}

	return content
}

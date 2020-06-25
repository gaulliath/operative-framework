package session

import (
	"encoding/base64"
	"encoding/json"
	"github.com/segmentio/ksuid"
	"time"
)

// Struct of module result
type OpfResults struct {
	Id         int              `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionId  int              `json:"-" gorm:"session_id"`
	ModuleName string           `json:"module_name"`
	ResultId   string           `json:"result_id" gorm:"primary_key:yes;column:result_id"`
	TargetId   string           `json:"target_id" gorm:"target_id"`
	Values     []OpfResultValue `json:"values"`
	ToJSON     json.RawMessage  `json:"to_json"`
	Notes      []Note
	Auxiliary  []string  `json:"-" sql:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

// Struct of module result value
type OpfResultValue struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Notes     []Note    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// Compact keys with separator
func (r *OpfResults) GetCompactKeys() string {
	compacted := ""
	for _, value := range r.Values {
		compacted = compacted + value.Key + base64.StdEncoding.EncodeToString([]byte(";operativeframework;"))[0:5]
	}
	return compacted
}

// Compact values with separator
func (r *OpfResults) GetCompactValues() string {
	compacted := ""
	for _, value := range r.Values {
		compacted = compacted + value.Value + base64.StdEncoding.EncodeToString([]byte(";operativeframework;"))[0:5]
	}
	return compacted
}

// Get result keys
func (r *OpfResults) GetKeys() []string {
	keys := []string{}
	for _, value := range r.Values {
		keys = append(keys, value.Key)
	}

	return keys
}

func (r *OpfResults) Set(key string, value string) {
	r.Values = append(r.Values, OpfResultValue{
		Key:       key,
		Value:     value,
		Notes:     []Note{},
		CreatedAt: time.Now(),
	})
}

func (r *OpfResults) Save(module Module, target *Target) bool {
	r.ModuleName = module.Name()
	if !target.ResultExist(*r) {
		target.Results[r.ModuleName] = append(target.Results[r.ModuleName], r)
		targets, err := target.Sess.FindLinked(r.ModuleName, *r)
		if err == nil {
			for _, id := range targets {
				target.Link(Linking{
					TargetId:       id,
					TargetResultId: r.ResultId,
				})
			}
		}

		i := target.Sess.GetCurrentInstance()
		i.SetResults(r.Values)
	}
	module.SetExport(*r)
	return true
}

func (target *Target) NewResult() *OpfResults {
	var result OpfResults
	result.ResultId = "R_" + ksuid.New().String()
	result.ToJSON = result.JSON()
	result.TargetId = target.GetId()
	result.SessionId = target.Sess.GetId()
	result.CreatedAt = time.Now()

	return &result
}

func (s *Session) GetResultsAfter(results []*OpfResults, afterTime time.Time) []*OpfResults {
	var scopes []*OpfResults
	for _, result := range results {
		if result.CreatedAt.After(afterTime) {
			scopes = append(scopes, result)
		}
	}

	return scopes
}

func (s *Session) GetResultsBefore(results []*OpfResults, beforeTime time.Time) []*OpfResults {
	var scopes []*OpfResults
	for _, result := range results {
		if result.CreatedAt.Before(beforeTime) {
			scopes = append(scopes, result)
		}
	}

	return scopes
}

func (r *OpfResults) Bytes() ([]byte, error) {
	content, err := json.Marshal(r.Values)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (r *OpfResults) JSON() json.RawMessage {
	content, err := json.Marshal(r.Values)
	if err != nil {
		return nil
	}
	return content
}

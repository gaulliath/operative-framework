package session

import (
	"github.com/segmentio/ksuid"
	"time"
)

type Instance struct {
	Id        string             `json:"id"`
	Module    string             `json:"module"`
	Results   [][]OpfResultValue `json:"results"`
	CreatedAt time.Time          `json:"created_at"`
}

// Get current session module execution instance
func (s *Session) GetCurrentInstance() *Instance {
	return s.CurrentInstance
}

// Get instance module name
func (i *Instance) GetModuleName() string {
	return i.Module
}

// Set result to instance
func (i *Instance) SetResults(r []OpfResultValue) {
	i.Results = append(i.Results, r)
}

// Create a new module execution instance
func (s *Session) NewInstance(module string) *Instance {
	i := &Instance{
		Id:        "INST_" + ksuid.New().String(),
		Module:    module,
		CreatedAt: time.Now(),
	}

	s.SetInstance(i)
	return i
}

// Set module execution instance to current
func (s *Session) SetInstance(i *Instance) *Session {
	s.CurrentInstance = i
	s.Instances = append(s.Instances, i)
	return s
}

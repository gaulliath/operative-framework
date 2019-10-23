package session

import (
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"strings"
	"time"
)

type Interval struct {
	Id               string        `json:"id"`
	S                *Session      `json:"-"`
	Commands         string        `json:"commands"`
	Delay            int           `json:"delay"`
	Time             time.Duration `json:"time"`
	ExecutionNumbers int           `json:"execution_numbers"`
	LastRun          time.Time     `json:"last_run"`
	NextRun          time.Time     `json:"next_run"`
	Activated        bool          `json:"activated"`
}

func (s *Session) NewInterval(command string) *Interval {
	newInterval := &Interval{
		Id:               "i_" + ksuid.New().String(),
		S:                s,
		ExecutionNumbers: 0,
		Time:             time.Minute,
		LastRun:          time.Now(),
		NextRun:          time.Now(),
		Activated:        false,
	}
	newInterval.SetCommand(command)
	s.Interval = append(s.Interval, newInterval)
	return newInterval
}

func (s *Session) WaitInterval() {
	for {
		time.Sleep(5 * time.Second)
		for _, interval := range s.Interval {
			if interval.Activated == true {
				interval.Up()
			}
		}
	}
}

func (s *Session) GetInterval(id string) (*Interval, error) {
	for _, interval := range s.Interval {
		if interval.Id == id {
			return interval, nil
		}
	}
	return nil, errors.New("This interval ID as not found.")
}

func (i *Interval) SetCommand(command string) *Interval {
	command = strings.TrimLeft(command, `"`)
	command = strings.TrimRight(command, `"`)
	i.Commands = command
	return i
}

func (i *Interval) GetCommand() string {
	return i.Commands
}

func (i *Interval) SetDelay(delay int) *Interval {
	i.Delay = delay
	return i
}

func (i *Interval) GetDelay() int {
	return i.Delay
}

func (i *Interval) Up() bool {
	timeNow := time.Now()
	if i.Activated == false {
		i.Activated = true
		i.NextRun = timeNow.Add(time.Duration(i.GetDelay()) * i.Time)
	} else {
		if timeNow.Equal(i.NextRun) || timeNow.After(i.NextRun) {
			i.S.Stream.Verbose = false
			if strings.Contains(i.GetCommand(), ";") {
				for _, command := range strings.Split(i.GetCommand(), ";") {
					i.S.ParseCommand(strings.TrimSpace(command))
				}
			} else {
				i.S.ParseCommand(strings.TrimSpace(i.GetCommand()))
			}
			i.S.Stream.Verbose = true
			i.LastRun = timeNow
			i.NextRun = timeNow.Add(time.Duration(i.GetDelay()) * i.Time)
		}
	}
	return false
}

func (i *Interval) Down() *Interval {
	i.Activated = false
	return i
}

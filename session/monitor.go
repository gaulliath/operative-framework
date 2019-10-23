package session

import (
	"errors"
	"github.com/graniet/go-pretty/table"
	"github.com/segmentio/ksuid"
	"os"
	"strings"
	"time"
)

type Monitors []*Monitor

type Monitor struct {
	Session   *Session
	MonitorId string
	Search    string
	Status    bool
	Result    []*TargetResults
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Session) NewMonitor(search string) *Monitor {

	newMonitor := Monitor{
		MonitorId: "M_" + ksuid.New().String(),
		Session:   s,
		Status:    false,
		Search:    search,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.Monitors = append(s.Monitors, &newMonitor)
	return &newMonitor
}

func (s *Session) GetMonitor(monitorId string) (*Monitor, error) {
	for _, monitor := range s.Monitors {
		if monitor.MonitorId == monitorId {
			return monitor, nil
		}
	}
	return nil, errors.New("Monitor '" + monitorId + "' not found in current session")
}

func (s *Session) GetMonitors() Monitors {
	return s.Monitors
}

func (s *Session) DeleteMonitor(monitorId string) {
	newMonitor := Monitors{}
	for _, monitor := range s.Monitors {
		if monitor.MonitorId != monitorId {
			newMonitor = append(newMonitor, monitor)
		}
	}
	s.Monitors = newMonitor
}

func (s *Session) WaitMonitor() {
	for {
		time.Sleep(10 * time.Second)
		for _, monitor := range s.Monitors {
			if monitor.Status == true {
				monitor.Checking()
			}
		}
	}
}

func (m *Monitor) Up() {
	m.Status = true
	m.UpdatedAt = time.Now()
}

func (m *Monitor) Down() {
	m.Status = false
	m.UpdatedAt = time.Now()
}

func (m *Monitor) ViewResults() {
	t := m.Session.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{40, 30, 30, 30})
	headerRow := table.Row{}
	for _, result := range m.Result {
		resRow := table.Row{}
		separator := m.Session.GetSeparator()
		header := strings.Split(result.Header, separator)
		res := strings.Split(result.Value, separator)
		if len(headerRow) < 1 {
			for _, h := range header {
				headerRow = append(headerRow, h)
			}
			headerRow = append(headerRow, "result_id")
			headerRow = append(headerRow, "target_id")
			t.AppendHeader(headerRow)
		}
		for _, r := range res {
			resRow = append(resRow, r)
		}
		resRow = append(resRow, result.ResultId)
		resRow = append(resRow, result.TargetId)
		t.AppendRow(resRow)
	}
	m.Session.Stream.Render(t)
}

func (m *Monitor) HasResult(resultId string) bool {
	for _, result := range m.Result {
		if result.ResultId == resultId {
			return true
		}
	}
	return false
}

func (m *Monitor) Checking() {
	for _, target := range m.Session.Targets {
		if len(target.Results) > 0 {
			for _, results := range target.Results {
				if len(results) > 0 {
					for _, result := range results {
						if result.CreatedAt.After(m.CreatedAt) {
							if strings.Contains(strings.ToLower(result.Value), strings.ToLower(m.Search)) {
								if !m.HasResult(result.ResultId) {
									m.Result = append(m.Result, result)
									m.UpdatedAt = time.Now()
									m.Session.NewEvent("MONITOR_MATCH", "Monitor '"+m.MonitorId+"' as found new matching with result '"+result.ResultId+"'")
								}
							}
						}
					}
				}
			}
		}
	}
}

package engine

import (
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"github.com/labstack/gommon/color"
	"os"
)

func CommandBase(line string, s *session.Session) bool{
	if line == "info session"{
		ViewInformation(s)
		return true
	}
	return false
}

func ViewInformation(s *session.Session){
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Name",
		"Value",
	})
	apiStatus := color.Red("offline")
	if s.Information.ApiStatus{
		apiStatus = color.Green("online")
	}
	t.AppendRow(table.Row{
		"API",
		apiStatus,
	})
	t.AppendRow(table.Row{
		"EVENT(S)",
		s.Information.Event,
	})
	t.AppendRow(table.Row{
		"MODULE(S)",
		len(s.Modules),
	})
	t.AppendRow(table.Row{
		"TARGET(S)",
		len(s.Targets),
	})
	s.Stream.Render(t)
}

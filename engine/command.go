package engine

import (
	"github.com/gorilla/mux"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/api"
	"github.com/graniet/operative-framework/session"
	"github.com/labstack/gommon/color"
	"os"
)

func CommandBase(line string, s *session.Session) bool{
	if line == "info session"{
		ViewInformation(s)
		return true
	} else if line== "info api"{
		ViewApiInformation(s)
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

func ViewApiInformation(s *session.Session){
	a := api.PushARestFul(s)
	r := a.LoadRouter()
	ta := s.Stream.GenerateTable()
	ta.SetOutputMirror(os.Stdout)
	ta.AppendHeader(table.Row{
		"Endpoint",
	})
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		ta.AppendRow(table.Row{
			t,
		})
		return nil
	})
	s.Stream.Render(ta)
}

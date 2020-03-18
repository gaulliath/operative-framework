package session

import (
	"errors"
	"github.com/graniet/go-pretty/table"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	"os"
	"path/filepath"
	"strings"
)

type WebHook struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Events []string `json:"events"`
	URL    string   `json:"url"`
	Method string   `json:"method"`
	Status bool     `json:"status"`
}

func (w *WebHook) GetId() string {
	return w.Id
}

func (w *WebHook) GetName() string {
	return w.Name
}

func (w *WebHook) GetEvents() string {
	return strings.Join(w.Events, ",")
}

func (w *WebHook) GetURL() string {
	return w.URL
}

func (w *WebHook) GetMethod() string {
	return w.Method
}

func (w *WebHook) GetStatus() bool {
	return w.Status
}

func (w *WebHook) SetStatus(s bool) *WebHook {
	w.Status = s
	return w
}

func (s *Session) PutWebHook(wh WebHook) {
	s.WebHooks = append(s.WebHooks, &wh)
}

func (w *WebHook) Up() {
	w.SetStatus(true)
}

func (w *WebHook) Down() {
	w.SetStatus(false)
}

func (s *Session) ListWebHooks() {
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"ID",
		"NAME",
		"METHOD",
		"EVENTS",
		"STATUS",
	})
	for _, webhook := range s.WebHooks {
		t.AppendRow(table.Row{
			webhook.GetId(),
			webhook.GetName(),
			webhook.GetEvents(),
			webhook.GetMethod(),
			webhook.GetStatus(),
		})
	}
	s.Stream.Render(t)
}

func (s *Session) AutoStartWebHook() {
	for _, wh := range s.WebHooks {
		wh.Up()
	}
}

func (s *Session) AutoStopWebHook() {
	for _, wh := range s.WebHooks {
		wh.Down()
	}
}

func (s *Session) GetWebHook(id string) (*WebHook, error) {
	for _, wh := range s.WebHooks {
		if wh.GetId() == id {
			return wh, nil
		}
	}
	return &WebHook{}, errors.New("Webhook as been not found.")
}

func (s *Session) GetWebHookByName(name string) (*WebHook, error) {
	for _, wh := range s.WebHooks {
		if wh.GetName() == name {
			return wh, nil
		}
	}
	return &WebHook{}, errors.New("Webhook '" + name + "' as been not found.")
}

func (s *Session) LoadWebHook() {

	validators := []string{
		"NAME",
		"EVENTS",
		"URL",
	}
	var events []string
	method := "POST"

	matches, _ := filepath.Glob(s.Config.Common.BaseDirectory + "webhooks/*.conf")
	for _, match := range matches {
		configuration, err := godotenv.Read(match)
		if err == nil {
			for _, validate := range validators {
				if _, ok := configuration[validate]; !ok {
					s.Stream.Error("Element '" + validate + "' as not found in webhook: '" + match + "'")
					return
				}
			}

			if strings.Contains(configuration["EVENTS"], ",") {
				events = strings.Split(configuration["EVENTS"], ",")
			} else {
				events = append(events, configuration["EVENTS"])
			}

			if _, ok := configuration["METHOD"]; ok {
				method = configuration["METHOD"]
			}

			s.PutWebHook(WebHook{
				Id:     "WH_" + ksuid.New().String(),
				Name:   configuration["NAME"],
				Events: events,
				URL:    configuration["URL"],
				Method: method,
				Status: false,
			})
		}
	}
}

func (s *Session) SendToWebHook(w *WebHook, event *Event) error {
	client := GetOpfClient()
	client.SetUserAgent("Operative-Framework : " + s.Version)
	_, err := client.SetData(event.JSON)
	if err == nil {
		_, err = client.Perform("POST", w.GetURL())
		return err
	}
	return err
}

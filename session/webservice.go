package session

import (
	"errors"
	"github.com/graniet/go-pretty/table"
	"github.com/joho/godotenv"
	"os"
	"os/user"
	"path/filepath"
)

type WebService struct {
	Name 		string `json:"name"`
	URL    		string `json:"url"`
}

func (s *Session) ParseWebServiceConfig() {
	u, _ := user.Current()
	matches, _ := filepath.Glob(u.HomeDir + "/.opf/external/webservices/*.conf")
	for _, match := range matches {
		configuration, err := godotenv.Read(match)
		if err == nil {
			s.WebServices = append(s.WebServices, WebService{
				Name:      configuration["NAME"],
				URL:       configuration["URL"],
			})
		}
	}
}

func (s *Session) GetWebService(name string) (WebService, error){
	for _, webservice := range s.WebServices {
		if webservice.Name == name {
			return webservice, nil
		}
	}
	return WebService{}, errors.New("Web service '"+name+"' as not found")
}

func (s *Session) ListWebServices() {
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"NAME",
		"URL",
	})

	for _, ws := range s.WebServices {
		t.AppendRow(table.Row{
			ws.Name,
			ws.URL,
		})
	}

	s.Stream.Render(t)
}

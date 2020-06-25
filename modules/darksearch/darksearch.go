package darksearch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
)

type DarkSearchModule struct {
	session.SessionModule
	Session *session.Session `json:"-"`
	Stream  *session.Stream  `json:"-"`
}

type DarkSearchResults struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	Data        []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Description string `json:"description"`
	} `json:"data"`
}

func PushDarkSearchModule(s *session.Session) *DarkSearchModule {
	mod := DarkSearchModule{
		Session: s,
		Stream:  &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Text to search", "", true, session.STRING)
	return &mod
}

func (module *DarkSearchModule) Name() string {
	return "dark_search"
}

func (module *DarkSearchModule) Description() string {
	return "Retrieve a results from TOR hidden service with DarkSearch.io"
}

func (module *DarkSearchModule) Author() string {
	return "Tristan Granier"
}

func (module *DarkSearchModule) GetType() []string {
	return []string{
		session.T_TARGET_TEXT,
	}
}

func (module *DarkSearchModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *DarkSearchModule) Start() {
	target, err := module.GetParameter("TARGET")
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	text, err := module.Session.GetTarget(target.Value)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	u := "https://darksearch.io/api/search?query=" + url.QueryEscape(text.GetName())
	client := http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	res, err := client.Do(req)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	if res.StatusCode != 200 {
		module.Stream.Error("Status Code: " + strconv.Itoa(res.StatusCode))
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	var results DarkSearchResults

	err = json.Unmarshal(body, &results)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Title",
		"Resume",
		"Link",
	})
	t.SetAllowedColumnLengths([]int{30, 30, 30})
	for _, element := range results.Data {
		t.AppendRow(table.Row{
			element.Title,
			element.Description,
			element.Link,
		})

		result := text.NewResult()
		result.Set("title", element.Title)
		result.Set("Resume", element.Description)
		result.Set("Link", element.Link)
		result.Save(module, text)
	}

	module.Stream.Render(t)
	return
}

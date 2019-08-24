package google

import (
	"net/http"
	u "net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
)

type GoogleTwitterModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushGoogleTwitterModule(s *session.Session) *GoogleTwitterModule {
	mod := GoogleTwitterModule{
		sess:   s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Search text", "", true, session.STRING)
	mod.CreateNewParam("LIMIT", "Limit search", "10", false, session.STRING)
	return &mod
}

func (module *GoogleTwitterModule) Name() string {
	return "google.twitter"
}

func (module *GoogleTwitterModule) Author() string {
	return "Tristan Granier"
}

func (module *GoogleTwitterModule) Description() string {
	return "Find result from google search engine"
}

func (module *GoogleTwitterModule) GetType() string {
	return "person"
}

func (module *GoogleTwitterModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *GoogleTwitterModule) Start() {
	paramEnterprise, _ := module.GetParameter("TARGET")
	target, err := module.sess.GetTarget(paramEnterprise.Value)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	if target.GetType() != module.GetType() {
		module.Stream.Error("Target with type '" + target.GetType() + "' isn't valid module need '" + module.GetType() + "' type.")
		return
	}

	paramLimit, _ := module.GetParameter("LIMIT")
	url := "https://www.google.com/search?num=" + paramLimit.Value + "&start=0&hl=en&q=" + u.QueryEscape(target.GetName()+" site:twitter.com")
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		module.Stream.Error("A error as been occurred with a target.")
		return
	}

	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{60, 100})
	t.AppendHeader(table.Row{"Link", "Username"})
	resultFound := 0
	doc.Find(".g").Each(func(i int, s *goquery.Selection) {
		link := s.Find("cite").Text()
		if strings.Contains(link, "twitter") {
			var username string
			if strings.Contains(link, "://") {
				username = strings.Split(link, "://")[1]
			} else {
				username = link
			}
			username = strings.Split(username, "/")[1]
			separator := target.GetSeparator()
			t.AppendRow([]interface{}{link, username})
			result := session.TargetResults{
				Header: "LINK" + separator + "USERNAME",
				Value:  link + separator + username,
			}
			target.Save(module, result)
			resultFound = resultFound + 1
			module.Results = append(module.Results, link)
		}
	})
	module.Stream.Render(t)
	module.Stream.Success(strconv.Itoa(resultFound) + " result(s) found with a google search engine.")
}

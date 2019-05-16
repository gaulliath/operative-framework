package google_search

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/go-pretty/table"
	"net/http"
	u "net/url"
	"os"
	"strconv"
	"strings"
)

type GoogleSearchModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushGoogleSearchModule(s *session.Session) *GoogleSearchModule{
	mod := GoogleSearchModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Search text", "", true, session.STRING)
	mod.CreateNewParam("LIMIT", "Limit search", "10", false, session.STRING)
	return &mod
}


func (module *GoogleSearchModule) Name() string{
	return "google_search"
}

func (module *GoogleSearchModule) Author() string{
	return "Tristan Granier"
}

func (module *GoogleSearchModule) Description() string{
	return "Find result from google search engine"
}

func (module *GoogleSearchModule) GetType() string{
	return "text"
}

func (module *GoogleSearchModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}


func (module *GoogleSearchModule) Start(){
	paramEnterprise, _ := module.GetParameter("TARGET")
	target, err := module.sess.GetTarget(paramEnterprise.Value)
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}

	if target.GetType() != module.GetType(){
		module.Stream.Error("Target with type '"+target.GetType()+"' isn't valid module need '"+module.GetType()+"' type.")
		return
	}

	paramLimit, _ := module.GetParameter("LIMIT")
	url := "https://encrypted.google.com/search?num=" + paramLimit.Value + "&start=0&hl=en&q=" + u.QueryEscape(target.GetName())
	res, err := http.Get(url)
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
	t.SetAllowedColumnLengths([]int{30, 100,})
	t.AppendHeader(table.Row{"Name", "Link"})

	resultFound := 0
	doc.Find(".g").Each(func(i int, s *goquery.Selection) {
		line := s.Find("h3").Text()
		name := strings.TrimSpace(line)
		link := s.Find("cite").Text()
		separator := target.GetSeparator()
		t.AppendRow([]interface{}{name, link})
		result := session.TargetResults{
			Header: "Name" + separator + "Link",
			Value: name + separator + link,
		}
		target.Save(module, result)
		resultFound = resultFound + 1
		module.Results = append(module.Results, link)
	})
	module.Stream.Render(t)
	module.Stream.Success(strconv.Itoa(resultFound) + " result(s) found with a google search engine.")
}
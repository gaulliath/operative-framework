package societe_com

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"net/http"
	"os"
	"strings"
)

type SocieteComModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushSocieteComModuleModule(s *session.Session) *SocieteComModule{
	mod := SocieteComModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Person name eg: Jhon Doe", "", true, session.STRING)
	mod.CreateNewParam("limit", "Limit search", "10", false, session.STRING)
	return &mod
}

func (module *SocieteComModule) Name() string{
	return "societe_com"
}

func (module *SocieteComModule) Description() string{
	return "Search possible enterprise on french network"
}

func (module *SocieteComModule) Author() string{
	return "Tristan Granier"
}

func (module *SocieteComModule) GetType() string{
	return "person"
}

func (module *SocieteComModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *SocieteComModule) Start(){
	paramPerson, _ := module.GetParameter("TARGET")
	target, err := module.sess.GetTarget(paramPerson.Value)
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}

	if target.GetType() != module.GetType(){
		module.Stream.Error("Target with type '"+target.GetType()+"' isn't valid module need '"+module.GetType()+"' type.")
		return
	}

	person := strings.Replace(target.Name," ", "%20", -1)

	paramLimit, _ := module.GetParameter("limit")
	url := "https://www.google.com/search?num=" + paramLimit.Value + "&start=0&hl=en&q=intext%3A\"" + person +"\"+site%3Asociete.com+intext%3A\"siret\"+"
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
	t.AppendHeader(table.Row{
		"ENTERPRISE",
	})
	doc.Find(".g").Each(func(i int, s *goquery.Selection) {
		line := s.Find("h3").Text()
		line = strings.Split(line, "(")[0]
		line = strings.TrimSpace(line)

		result := session.TargetResults{
			Header: "enterprise",
			Value: line,
		}
		module.Results = append(module.Results, line)
		target.Save(module, result)

		t.AppendRow(table.Row{
			line,
		})
	})
	module.Stream.Render(t)
}
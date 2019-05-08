package pastebin

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/go-pretty/table"
	"net/http"
	"os"
	"strings"
)

type PasteBin struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushPasteBinModule(s *session.Session) *PasteBin{
	mod := PasteBin{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Email address", "", true, session.STRING)
	mod.CreateNewParam("limit", "Limit search", "10", false, session.STRING)
	return &mod
}

func (module *PasteBin) Name() string{
	return "pastebin"
}

func (module *PasteBin) Description() string{
	return "Check possible text on pastebin.com"
}

func (module *PasteBin) Author() string{
	return "Tristan Granier"
}

func (module *PasteBin) GetType() string{
	return "text"
}


func (module *PasteBin) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *PasteBin) Start(){
	paramEmail, _ := module.GetParameter("TARGET")
	target, err := module.sess.GetTarget(paramEmail.Value)
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}

	if target.GetType() != module.GetType(){
		module.Stream.Error("Target with type '"+target.GetType()+"' isn't valid module need '"+module.GetType()+"' type.")
		return
	}

	paramLimit, _ := module.GetParameter("limit")
	urlEnd := strings.Replace(target.GetName(), "@", "%40", -1)
	urlEnd = strings.Replace(urlEnd, " ", "%20", -1)
	url := "https://encrypted.google.com/search?num=" + paramLimit.Value + "&start=0&hl=en&q=site%3Apastebin.com%20\""+urlEnd+"\""
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
	t.SetAllowedColumnLengths([]int{60,})
	t.AppendHeader(table.Row{"Link"})

	resultFound := 0
	doc.Find(".g").Each(func(i int, s *goquery.Selection) {
		link := strings.TrimSpace(s.Find("cite").Text())
		separator := target.GetSeparator()
		t.AppendRow(table.Row{
			link,
		})
		result := session.TargetResults{
			Header: "link" + separator,
			Value: link + separator,
		}
		module.Results = append(module.Results, link)
		target.Save(module, result)
		resultFound = resultFound + 1
	})
	module.Stream.Render(t)
}

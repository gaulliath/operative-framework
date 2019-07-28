package pastebin_email

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/go-pretty/table"
	"net/http"
	"os"
	"strings"
)

type PasteBinEmail struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushPasteBinEmailModule(s *session.Session) *PasteBinEmail{
	mod := PasteBinEmail{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Email address", "", true, session.STRING)
	mod.CreateNewParam("limit", "Limit search", "10", false, session.STRING)
	return &mod
}

func (module *PasteBinEmail) Name() string{
	return "pastebin_email"
}

func (module *PasteBinEmail) Description() string{
	return "Check possible email on pastebin.com"
}

func (module *PasteBinEmail) Author() string{
	return "Tristan Granier"
}

func (module *PasteBinEmail) GetType() string{
	return "email"
}


func (module *PasteBinEmail) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *PasteBinEmail) Start(){
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
	url := "https://www.google.com/search?num=" + paramLimit.Value + "&start=0&hl=en&q=site%3Apastebin.com%20\""+urlEnd+"\""
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
	t.SetAllowedColumnLengths([]int{60,})
	t.AppendHeader(table.Row{"Link"})

	resultFound := 0
	doc.Find(".g").Each(func(i int, s *goquery.Selection) {
		link := strings.TrimSpace(s.Find("cite").Text())
		t.AppendRow(table.Row{
			link,
		})
		result := session.TargetResults{
			Header: "link",
			Value: link,
		}
		module.Results = append(module.Results, link)
		target.Save(module, result)
		resultFound = resultFound + 1
	})
	module.Stream.Render(t)
}

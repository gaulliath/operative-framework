package pastebin

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
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
	urlEnd = strings.Replace(urlEnd, " ", "+", -1)
	url := "https://www.google.com/search?q=site%3Apastebin.com+\""+urlEnd+"\"&num="+paramLimit.Value+"&hl=com"
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
	sel := doc.Find("div.g")
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		link = strings.Trim(link, " ")
		if link != "" && link != "#" {
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
		}
	}
	module.Stream.Render(t)
}

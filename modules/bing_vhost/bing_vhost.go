package bing_vhost

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"net/http"
	"os"
	"strings"
)

type BingVirtualHostModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushBingVirtualHostModule(s *session.Session) *BingVirtualHostModule{
	mod := BingVirtualHostModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Target argument for vhost checking", "", true, session.STRING)
	return &mod
}


func (module *BingVirtualHostModule) Name() string{
	return "bing_vhost"
}

func (module *BingVirtualHostModule) Author() string{
	return "Tristan Granier"
}

func (module *BingVirtualHostModule) Description() string{
	return "Checking possible virtual host with target IP"
}

func (module *BingVirtualHostModule) GetType() string{
	return "ip_address"
}

func (module *BingVirtualHostModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *BingVirtualHostModule) Start(){
	ipAddress, err := module.GetParameter("TARGET")
	if err != nil{
		module.Stream.Error("Argument 'TARGET' isn't valid.")
		return
	}

	target, err := module.sess.GetTarget(ipAddress.Value)
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}

	if target.GetType() != module.GetType(){
		module.Stream.Error("Target with type '"+target.GetType()+"' isn't valid module need '"+module.GetType()+"' type.")
		return
	}

	url := "https://www.bing.com/search?q=ip%3a" + target.GetName()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.3")
	res, err := client.Do(req)
	if err != nil {
		module.Stream.Error("'URL' can't be reached.")
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		module.Stream.Error("'URL' can't be reached.")
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		module.Stream.Error("A error as been occurred with a target.")
		return
	}

	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Link"})
	doc.Find("cite").Each(func(i int, s *goquery.Selection) {
		line := strings.TrimSpace(s.Text())
		t.AppendRow(table.Row{line})
		result := session.TargetResults{
			Header: "Link",
			Value: line,
		}
		target.Save(module, result)
		module.Results = append(module.Results, line)
	})
	module.Stream.Render(t)

}

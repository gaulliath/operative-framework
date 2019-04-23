package metatag_spider

import (
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/go-pretty/table"
	"net/http"
	"net/url"
	"os"
	"github.com/PuerkitoBio/goquery"
)

type MetaTagModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushMetaTagModule(s *session.Session) *MetaTagModule{
	mod := MetaTagModule{
		sess: s,
		Stream: &s.Stream,
	}
	mod.CreateNewParam("TARGET", "Target website URL","", true, session.STRING)
	return &mod
}

func (module *MetaTagModule) Name() string{
	return "metatag_spider"
}

func (module *MetaTagModule) Description() string{
	return "Crawl a meta tags elements for selected target"
}

func (module *MetaTagModule) Author() string{
	return "Tristan Granier"
}

func (module *MetaTagModule) GetType() string{
	return "url"
}

func (module *MetaTagModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *MetaTagModule) Start(){
	targetUrl, err := module.GetParameter("TARGET")
	if err != nil{
		module.Stream.Error("Argument 'URL' can't be parsed.")
		return
	}

	target, err := module.sess.GetTarget(targetUrl.Value)
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}

	if target.GetType() != module.GetType(){
		module.Stream.Error("Target with type '"+target.GetType()+"' isn't valid module need '"+module.GetType()+"' type.")
		return
	}

	_, err = url.ParseRequestURI(target.GetName())
	if err != nil {
		module.Stream.Error("Argument 'URL' isn't valid.")
		module.SetParameter("TARGET", "")
		return
	}

	res, err := http.Get(target.GetName())
	if err != nil{
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}

	tagFound := make(map[string]string)

	doc.Find("meta").Each(func(i int, s *goquery.Selection){
		tagContent,_ := s.Attr("content")
		tagName,_ := s.Attr("name")
		separator := target.GetSeparator()
		if _, ok := tagFound[tagName]; !ok{
			if tagName != "" {
				tagFound[tagName] = tagContent
				result := session.TargetResults{
					Header: "KEY" + separator + "VALUE",
					Value: tagName + separator + tagContent,
				}
				target.Save(module, result)
			}
		}
	})
	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{0, 60,})
	t.AppendHeader(table.Row{"KEY", "VALUE"})

	if len(tagFound) > 0 {
		for name, content := range tagFound {
			t.AppendRow([]interface{}{name, content})
		}
	} else {
		t.AppendRow([]interface{}{"No meta tag found", "No meta tag found"})
	}
	module.Stream.Render(t)
}
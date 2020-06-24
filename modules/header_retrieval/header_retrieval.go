package header_retrieval

import (
	"net/url"
	"os"
	"strings"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"github.com/imroc/req"
)

type HeaderRetrievalModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushModuleHeaderRetrieval(s *session.Session) *HeaderRetrievalModule {
	mod := HeaderRetrievalModule{
		sess:   s,
		Stream: &s.Stream,
	}
	mod.CreateNewParam("TARGET", "target URL", "", true, session.STRING)
	mod.CreateNewParam("METHOD", "Method used", "GET", true, session.STRING)
	return &mod
}

func (module *HeaderRetrievalModule) Name() string {
	return "get.header"
}

func (module *HeaderRetrievalModule) Author() string {
	return "Tristan Granier"
}

func (module *HeaderRetrievalModule) Description() string {
	return "Get headers from selected URL"
}

func (module *HeaderRetrievalModule) GetType() []string {
	return []string{
		session.T_TARGET_URL,
	}
}

func (module *HeaderRetrievalModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *HeaderRetrievalModule) Start() {
	paramURL, _ := module.GetParameter("TARGET")
	target, err := module.sess.GetTarget(paramURL.Value)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	_, err = url.ParseRequestURI(target.GetName())
	if err != nil {
		module.Stream.Error("Argument 'TARGET' isn't valid.")
		_, _ = module.SetParameter("TARGET", "")
		return
	}

	header := req.Header{
		"Accept":     "application/json",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
	}

	if !strings.Contains(target.GetName(), "://") {
		target.Name = "http://" + target.Name
	}

	r, err := req.Get(target.GetName(), header)
	if err != nil {
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}
	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{0, 60})
	t.AppendHeader(table.Row{"KEY", "VALUE"})
	for index, header := range r.Response().Header {
		t.AppendRow([]interface{}{index, header})
		if len(header) > 0 {
			result := target.NewResult()
			for _, l := range header {
				result.Set(index, l)
			}
			result.Save(module, target)
		}
	}
	module.Stream.Render(t)
}

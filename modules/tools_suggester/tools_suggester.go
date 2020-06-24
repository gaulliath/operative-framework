package tools_suggester

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
)

type ToolsSuggesterModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
	Lists  []WorldList      `json:"-"`
}

type WorldList struct {
	Name   string
	Path   string
	Tools  string
	Code   []int
	InText string
}

type Results struct {
	Url        string
	StatusCode int
	Tools      string
	Type       string
}

func PushModuleToolsSuggester(s *session.Session) *ToolsSuggesterModule {

	mod := ToolsSuggesterModule{
		sess:   s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "website target", "", true, session.STRING)
	mod.PushToLists("WordPress", "/wp-includes/", "wpscan", []int{200, 403}, "")
	mod.PushToLists("WordPress", "/wp-admin/", "wpscan", []int{200, 403}, "")
	mod.PushToLists("WordPress", "/readme.html", "wpscan", []int{200}, "WordPress")
	mod.PushToLists("Drupal", "/CHANGELOG.txt", "drupscan", []int{200}, "drupal")
	mod.PushToLists("Magento", "/frontend/default/", "Magescan", []int{200, 403}, "")
	mod.PushToLists("Magento", "/static/frontend/", "Magescan", []int{200, 403}, "")
	mod.PushToLists("Magento", "/magento_version", "Magescan", []int{200, 403}, "Magento")
	return &mod
}

func (module *ToolsSuggesterModule) PushToLists(name string, path string, tools string, code []int, intext string) {
	module.Lists = append(module.Lists, WorldList{
		Name:   name,
		Path:   path,
		Tools:  tools,
		Code:   code,
		InText: intext,
	})
	return
}

func (module *ToolsSuggesterModule) Name() string {
	return "tools.suggester"
}

func (module *ToolsSuggesterModule) Author() string {
	return "Tristan Granier"
}

func (module *ToolsSuggesterModule) Description() string {
	return "Search possible tools for exploit target"
}

func (module *ToolsSuggesterModule) GetType() []string {
	return []string{
		session.T_TARGET_WEBSITE,
	}
}

func (module *ToolsSuggesterModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *ToolsSuggesterModule) Start() {
	trg, err := module.GetParameter("TARGET")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	target, err := module.sess.GetTarget(trg.Value)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	if strings.Contains(target.GetName(), "://") {
		expProto := strings.Split(target.GetName(), "://")
		proto := expProto[0]
		expURL := ""
		if strings.Contains(target.GetName(), "/") {
			expURL = strings.Split(expProto[1], "/")[0]
			target.Name = proto + "://" + expURL
		}
	} else {

		if strings.Contains(target.GetName(), "/") {
			expURL := strings.Split(target.GetName(), "/")[0]
			target.Name = "https://" + expURL
		}
	}

	client := http.Client{}
	var results []Results
	for _, obj := range module.Lists {
		url := target.GetName() + obj.Path
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
		res, err := client.Do(req)
		if err != nil {
			continue
		}

		for _, code := range obj.Code {
			if res.StatusCode == code {
				if obj.InText != "" {
					body, err := ioutil.ReadAll(res.Body)
					if err != nil {
						continue
					}
					if strings.Contains(string(body), obj.InText) {
						results = append(results, Results{
							Url:        url,
							StatusCode: res.StatusCode,
							Tools:      obj.Tools,
							Type:       obj.Name,
						})
						trgRest := target.NewResult()
						trgRest.Set("URL", url)
						trgRest.Set("RESPONSE", strconv.Itoa(res.StatusCode))
						trgRest.Set("TOOLS", obj.Tools)
						trgRest.Set("CMS", obj.Name)
						trgRest.Save(module, target)

						module.Results = append(module.Results, obj.Tools)
					}
				} else {
					results = append(results, Results{
						Url:        url,
						StatusCode: res.StatusCode,
						Tools:      obj.Tools,
						Type:       obj.Name,
					})

					result := target.NewResult()
					result.Set("URL", url)
					result.Set("RESPONSE", strconv.Itoa(res.StatusCode))
					result.Set("TOOL", obj.Tools)
					result.Set("CMS", obj.Name)
					result.Save(module, target)
					module.Results = append(module.Results, obj.Tools)
					continue
				}
			}
		}
	}

	t := module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"URL",
		"RESPONSE",
		"TOOLS",
		"CMS",
	})
	if len(results) > 0 {
		for _, result := range results {
			t.AppendRow(table.Row{
				result.Url,
				result.StatusCode,
				result.Tools,
				result.Type,
			})
		}
		module.sess.Stream.Render(t)
	} else {
		module.sess.Stream.Warning("Not result found.")
	}
	return
}

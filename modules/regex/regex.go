package regex

import (
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"os"
	"regexp"
	"strings"
)

type FindWithRegexModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushFindWithRegexModule(s *session.Session) *FindWithRegexModule {
	mod := FindWithRegexModule{
		sess:   s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("search", "regex term e.g: ^[a-z]+\\[[0-9]+\\]$", "", true, session.STRING)
	mod.CreateNewParam("source", "or would you like to search?", "results", false, session.STRING)
	return &mod
}

func (module *FindWithRegexModule) Name() string {
	return "regex"
}

func (module *FindWithRegexModule) Description() string {
	return "Find by regex term in result(s) or target(s)"
}

func (module *FindWithRegexModule) Author() string {
	return "Tristan Granier"
}

func (module *FindWithRegexModule) GetType() []string {
	return []string{
		session.T_TARGET_COMMAND,
	}
}

func (module *FindWithRegexModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *FindWithRegexModule) Start() {

	term, err := module.GetParameter("SEARCH")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	searchIn, err := module.GetParameter("SOURCE")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	found := false

	sources := []string{
		"results",
		"targets",
	}

	for _, source := range sources {
		if strings.ToLower(source) == strings.ToLower(searchIn.Value) {
			found = true
		}
	}

	if found == false {
		module.Stream.Error("Source '" + searchIn.Value + "' as unknown please use: results or targets")
		return
	}

	s := module.sess

	targets := s.Targets
	var separator string
	var results []*session.OpfResults
	var r = regexp.MustCompile(term.Value)

	for _, target := range targets {
		separator = target.GetSeparator()
		modules := target.Results
		for _, moduleResults := range modules {
			for _, result := range moduleResults {
				keys := strings.Split(result.GetCompactKeys(), separator)
				values := strings.Split(result.GetCompactValues(), separator)
				found := false
				for _, key := range keys {
					if r.MatchString(key) {
						results = append(results, result)
						found = true
					}
				}

				for _, value := range values {
					if r.MatchString(value) {
						if found == false {
							results = append(results, result)
							found = true
						}
					}
				}
			}
		}
	}

	for _, result := range results {
		t := s.Stream.GenerateTable()
		t.SetOutputMirror(os.Stdout)
		t.SetAllowedColumnLengths([]int{30, 30})
		t.AppendRow(table.Row{
			"TARGET",
			result.TargetId,
		})
		t.AppendRow(table.Row{
			"RESULT",
			result.ResultId,
		})
		keys := strings.Split(result.GetCompactKeys(), separator)
		values := strings.Split(result.GetCompactValues(), separator)
		for index, key := range keys {
			realValue := ""
			if values[index] != "" {
				realValue = values[index]
			}
			t.AppendRow(table.Row{
				key,
				realValue,
			})
		}
		s.Stream.Render(t)
	}
	return
}

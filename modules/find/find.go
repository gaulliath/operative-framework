package find

import (
	"github.com/fatih/color"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"os"
	"strings"
)

type FindModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushFindModule(s *session.Session) *FindModule {
	mod := FindModule{
		sess:   s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("search", "search term e.g: Operative Framework", "", true, session.STRING)
	mod.CreateNewParam("source", "or would you like to search?", "results", false, session.STRING)
	return &mod
}

func (module *FindModule) Name() string {
	return "find"
}

func (module *FindModule) Description() string {
	return "Find by search term in result(s) or target(s)"
}

func (module *FindModule) Author() string {
	return "Tristan Granier"
}

func (module *FindModule) GetType() string {
	return "command"
}

func (module *FindModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *FindModule) Start() {

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
	var results []session.TargetResults

	for _, target := range targets {
		separator = target.GetSeparator()
		modules := target.Results
		for _, moduleResults := range modules {
			for _, result := range moduleResults {

				var temporary session.TargetResults
				temporary = *result

				keys := strings.Split(temporary.Header, separator)
				values := strings.Split(temporary.Value, separator)
				found := false

				for k1, key := range keys {
					if strings.Contains(key, term.Value) {
						key = strings.Replace(key, term.Value, color.RedString(term.Value), -1)
						keys[k1] = key
						temporary.Header = strings.Join(keys, separator)
						found = true
					}
				}

				for k2, value := range values {
					if strings.Contains(value, term.Value) {
						value = strings.Replace(value, term.Value, color.RedString(term.Value), -1)
						values[k2] = value
						temporary.Value = strings.Join(values, separator)
						found = true
					}
				}

				if found == true {
					results = append(results, temporary)
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
		keys := strings.Split(result.Header, separator)
		values := strings.Split(result.Value, separator)
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

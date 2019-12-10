package session

import (
	"errors"
	"github.com/graniet/go-pretty/table"
	"os"
)

func (s *Session) AddAlias(alias string, module string) {
	mod, err := s.SearchModule(module)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	lists := make(map[string]string)
	for a, als := range s.Alias {
		if als != mod.Name() {
			lists[a] = als
		}
	}

	lists[alias] = mod.Name()
	s.Alias = lists
	s.Stream.Success("Alias '" + alias + "' as created for module '" + mod.Name() + "'")
	return
}

func (s *Session) GetAlias(alias string) (string, error) {
	for a, m := range s.Alias {
		if a == alias {
			return m, nil
		}
	}
	return "", errors.New("Alias not found in current session.")
}

func (s *Session) ListAlias() {
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"ALIAS",
		"MODULE",
	})
	for a, m := range s.Alias {
		t.AppendRow(table.Row{
			a,
			m,
		})
	}
	s.Stream.Render(t)
	return
}

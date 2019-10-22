package session

import (
	"strings"
)

func (s *Session) ParseCommand(line string) []string {
	moduleName := strings.Split(line, " ")[0]
	module, errModule := s.SearchModule(moduleName)
	if errModule != nil {
		alias, err := s.GetAlias(moduleName)
		module, err = s.SearchModule(alias)
		if err != nil {
			if moduleName == "help" {
				module, err = s.SearchModule("session_help")
			}
		}
	}
	if strings.Contains(line, " ") {
		if strings.HasPrefix(line, "sh ") {
			LoadShCommandMenu(line, module, s)

		} else if strings.HasPrefix(line, "find ") || strings.HasPrefix(line, "regex ") {
			LoadFindCommandMenu(line, module, s)

		} else if strings.HasPrefix(line, "alias ") {
			LoadAliasMenu(line, module, s)

		} else if strings.HasPrefix(line, "note ") {
			LoadNoteMenu(line, module, s)

		} else if strings.HasPrefix(line, "target ") {
			LoadTargetMenu(line, module, s)

		} else if strings.HasPrefix(line, "interval ") {
			LoadIntervalCommandMenu(line, module, s)

		} else if strings.HasPrefix(line, "modules ") {
			LoadModuleByTypeMenu(line, module, s)

		} else if strings.HasPrefix(line, "analytics ") {
			LoadAnalyticsWebBased(line, module, s)
		} else {
			if errModule == nil {
				LoadModuleMenu(line, module, s)
			}
		}
	}
	if line == "events" {
		LoadEventsMenu(line, module, s)
		return nil
	}
	if moduleName == "help" {
		module.Start()
		filter, err := module.GetParameter("FILTER")
		if err == nil && filter.Value != "" {
			flt, err := s.SearchFilter(filter.Value)
			if err != nil {
				s.Stream.Error("Filter '" + filter.Value + "' as not found.")
				return nil
			}
			s.Stream.Success("Start filter '" + filter.Value + "'...")
			flt.Start(module)
		}
	}
	return nil
}

package session

import (
	"github.com/chzyer/readline"
	"os"
)

func (s *Session) ReadLineAutoCompleteType() func(string) []string{
	return func(line string) []string{
		return s.ListType()
	}
}

func (s *Session) ReadLineAutoCompleteListModules() func(string) []string{
	return func(line string) []string{
		modules := make([]string, 0)
		for _, module := range s.Modules{
			modules = append(modules, module.Name())
		}
		return modules
	}
}

func (s *Session) ReadLineAutoCompleteTargets() func(string) []string{
	return func(line string) []string{
		targets := make([]string, 0)
		for _, target := range s.Targets{
			targets = append(targets, target.TargetId)
		}
		return targets
	}
}

func (s *Session) PushPrompt(){
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("target",
			readline.PcItem("add",
				readline.PcItemDynamic(s.ReadLineAutoCompleteType())),
			readline.PcItem("delete"),
			readline.PcItem("update"),
			readline.PcItem("list"),
			readline.PcItem("links",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("modules"),
		),
		readline.PcItemDynamic(s.ReadLineAutoCompleteListModules(),
			readline.PcItem("list"),
			readline.PcItem("target",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("set",
				readline.PcItem("TARGET")),
			readline.PcItem("run"),
		),
		readline.PcItem("help"),
		readline.PcItem("info"),
	)
	s.Prompt = &readline.Config{
		Prompt:          "\033[90m[OPF v"+s.Version+"]:\033[0m ",
		HistoryFile:     os.Getenv("OPERATIVE_HISTORY"),
		InterruptPrompt: "^C",
		AutoComplete: completer,
		EOFPrompt:       "exit",
		HistorySearchFold:   true,
	}
}

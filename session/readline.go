package session

import (
	"github.com/chzyer/readline"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func (s *Session) ClearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin", "freebsd", "dragonfly", "netbsd", "openbsd", "solaris":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	case "windows":
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	default:
		// do nothing
	}
}

func (s *Session) ReadLineAutoCompleteType() func(string) []string {
	return func(line string) []string {
		return s.ListType()
	}
}

func (s *Session) ReadLineAutoCompleteListModules() func(string) []string {
	return func(line string) []string {
		modules := make([]string, 0)
		for _, module := range s.Modules {
			modules = append(modules, module.Name())
		}
		return modules
	}
}

func (s *Session) ReadLineAutoCompleteFilters() func(string) []string {
	return func(line string) []string {
		filters := make([]string, 0)
		for _, filter := range s.Filters {
			filters = append(filters, filter.Name())
		}
		return filters
	}
}

func (s *Session) ReadLineAutoCompleteTracker() func(string) []string {
	return func(line string) []string {
		trackers := make([]string, 0)
		for _, track := range s.Tracker.Tracked {
			trackers = append(trackers, track.Id)
		}
		return trackers
	}
}

func (s *Session) ReadLineAutoCompleteTargets() func(string) []string {
	return func(line string) []string {
		targets := make([]string, 0)
		for _, target := range s.Targets {
			targets = append(targets, target.TargetId)
		}
		return targets
	}
}

func (s *Session) ReadLineAutoCompleteInterval() func(string) []string {
	return func(line string) []string {
		intervals := make([]string, 0)
		for _, interval := range s.Interval {
			intervals = append(intervals, interval.Id)
		}
		return intervals
	}
}

func (s *Session) ReadLineAutoCompleteMonitor() func(string) []string {
	return func(line string) []string {
		monitors := make([]string, 0)
		for _, monitor := range s.Monitors {
			monitors = append(monitors, monitor.MonitorId)
		}
		return monitors
	}
}

func (s *Session) ReadLineAutoCompleteModuleResults() func(string) []string {
	return func(line string) []string {
		value := strings.Split(line, " ")
		var returnResult []string
		if len(value) < 4 {
			return []string{}
		}
		target, err := s.GetTarget(value[3])
		if err != nil {
			return []string{}
		}
		for name := range target.Results {
			returnResult = append(returnResult, name)
		}
		return returnResult
	}
}

func (s *Session) ReadLineAutoCompleteListAlias() func(string) []string {
	return func(line string) []string {
		var returnResult []string
		for name := range s.Alias {
			returnResult = append(returnResult, name)
		}
		return returnResult
	}
}

func (s *Session) ReadLineAutoCompleteCacheName() func(string) []string {
	return func(line string) []string {
		var cacheName []string
		file, _ := filepath.Glob(s.Config.Common.BaseDirectory + "cache/*")
		for _, name := range file {
			cacheName = append(cacheName, name)
		}
		return cacheName
	}
}

func (s *Session) ReadLineAutoCompleteListWebHooks() func(string) []string {
	return func(line string) []string {
		var returnResult []string
		for _, wh := range s.WebHooks {
			returnResult = append(returnResult, wh.GetId())
		}
		return returnResult
	}
}

func (s *Session) ReadLineAutoCompleteResults() func(string) []string {
	return func(line string) []string {
		value := strings.Split(line, " ")
		var returnResult []string
		if len(value) < 4 {
			if value[1] == "add" {
				for _, target := range s.Targets {
					for _, module := range target.Results {
						for _, result := range module {
							if len(result.Values) > 0 {
								returnResult = append(returnResult, result.ResultId)
							}
						}
					}
				}
				return returnResult
			}
			return []string{}
		}
		target, err := s.GetTarget(value[3])
		if err != nil {
			return []string{}
		}
		for _, module := range target.Results {
			for _, result := range module {
				if len(result.Values) > 0 {
					returnResult = append(returnResult, result.ResultId)
				}
			}
		}
		return returnResult
	}
}

func (s *Session) PushPrompt() {
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("modules"),
		readline.PcItem("webhooks"),
		readline.PcItem("webhook",
			readline.PcItem("up",
				readline.PcItemDynamic(s.ReadLineAutoCompleteListWebHooks())),
			readline.PcItem("down",
				readline.PcItemDynamic(s.ReadLineAutoCompleteListWebHooks()))),
		readline.PcItem("result",
			readline.PcItem("delete",
				readline.PcItemDynamic(s.ReadLineAutoCompleteResults()))),
		readline.PcItem("events"),
		readline.PcItem("monitor",
			readline.PcItem("generate"),
			readline.PcItem("list"),
			readline.PcItem("results",
				readline.PcItemDynamic(s.ReadLineAutoCompleteMonitor())),
			readline.PcItem("up",
				readline.PcItemDynamic(s.ReadLineAutoCompleteMonitor())),
			readline.PcItem("down",
				readline.PcItemDynamic(s.ReadLineAutoCompleteMonitor()))),
		readline.PcItem("analytics",
			readline.PcItem("up"),
			readline.PcItem("down")),
		readline.PcItem("interval",
			readline.PcItem("list"),
			readline.PcItem("generate"),
			readline.PcItem("set",
				readline.PcItemDynamic(s.ReadLineAutoCompleteInterval())),
			readline.PcItem("up",
				readline.PcItemDynamic(s.ReadLineAutoCompleteInterval())),
			readline.PcItem("down",
				readline.PcItemDynamic(s.ReadLineAutoCompleteInterval()))),
		readline.PcItem("alias",
			readline.PcItem("list"),
			readline.PcItem("add",
				readline.PcItemDynamic(s.ReadLineAutoCompleteListModules()))),
		readline.PcItem("note",
			readline.PcItem("add",
				readline.PcItemDynamic(s.ReadLineAutoCompleteResults())),
			readline.PcItem("view",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets()))),
		readline.PcItem("target",
			readline.PcItem("add",
				readline.PcItemDynamic(s.ReadLineAutoCompleteType())),
			readline.PcItem("delete",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("update",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("list"),
			readline.PcItem("type"),
			readline.PcItem("link",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets(),
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets()))),
			readline.PcItem("links",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("modules", readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("tag",
				readline.PcItem("add",
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
				readline.PcItem("list",
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets()))),
			readline.PcItem("view",
				readline.PcItem("result",
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets(),
						readline.PcItemDynamic(s.ReadLineAutoCompleteResults()))),
				readline.PcItem("results",
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets(),
						readline.PcItemDynamic(s.ReadLineAutoCompleteModuleResults()))),
				readline.PcItem("notes")),
		),
		readline.PcItemDynamic(s.ReadLineAutoCompleteListAlias(),
			readline.PcItem("list"),
			readline.PcItem("target",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("filter",
				readline.PcItemDynamic(s.ReadLineAutoCompleteFilters())),
			readline.PcItem("set",
				readline.PcItem("TARGET"),
				readline.PcItem("DISABLE_OUTPUT")),
			readline.PcItem("run"),
		),
		readline.PcItemDynamic(s.ReadLineAutoCompleteListModules(),
			readline.PcItem("list"),
			readline.PcItem("target",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("filter",
				readline.PcItemDynamic(s.ReadLineAutoCompleteFilters())),
			readline.PcItem("set",
				readline.PcItem("TARGET"),
				readline.PcItem("DISABLE_OUTPUT")),
			readline.PcItem("run"),
		),
		readline.PcItem("help"),
		readline.PcItem("env"),
		readline.PcItem("save"),
		readline.PcItem("load",
			readline.PcItemDynamic(s.ReadLineAutoCompleteCacheName())),
		readline.PcItem("info",
			readline.PcItem("session"),
			readline.PcItem("api")),
		readline.PcItem("api",
			readline.PcItem("run"),
			readline.PcItem("stop")),
		readline.PcItem("tracker",
			readline.PcItem("run"),
			readline.PcItem("stop"),
			readline.PcItem("list"),
			readline.PcItem("positions",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTracker())),
			readline.PcItem("select",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTracker())),
		),
	)
	s.Prompt = &readline.Config{
		Prompt:            "\033[90m[OPF v" + s.Version + "]:\033[0m ",
		HistoryFile:       s.Config.Common.HistoryFile,
		InterruptPrompt:   "^C",
		AutoComplete:      completer,
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	}
}

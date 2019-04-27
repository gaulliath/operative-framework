package session

import (
	"github.com/chzyer/readline"
	"github.com/graniet/go-pretty/table"
	"os"
	"strings"
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

func (s *Session) ReadLineAutoCompleteFilters() func(string) []string{
	return func(line string) []string{
		filters := make([]string, 0)
		for _, filter := range s.Filters{
			filters = append(filters, filter.Name())
		}
		return filters
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

func (s *Session) ReadLineAutoCompleteModuleResults() func(string) []string{
	return func(line string) []string{
		value := strings.Split(line, " ")
		var returnResult []string
		if len(value) < 4{
			return []string{
			}
		}
		target, err := s.GetTarget(value[3])
		if err != nil{
			return []string{
			}
		}
		for name := range target.Results{
			returnResult = append(returnResult, name)
		}
		return returnResult
	}
}

func (s *Session) ReadLineAutoCompleteResults() func(string) []string{
	return func(line string) []string{
		value := strings.Split(line, " ")
		var returnResult []string
		if len(value) < 4{
			return []string{
			}
		}
		target, err := s.GetTarget(value[3])
		if err != nil{
			return []string{
			}
		}
		for _, module := range target.Results{
			for _, result := range module{
				if result.Value != "" {
					returnResult = append(returnResult, result.ResultId)
				}
			}
		}
		return returnResult
	}
}

func (s *Session) PushPrompt(){
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("target",
			readline.PcItem("add",
				readline.PcItemDynamic(s.ReadLineAutoCompleteType())),
			readline.PcItem("delete",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("update",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("list"),
			readline.PcItem("links",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("modules",readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("view",
				readline.PcItem("result",
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets(),
						readline.PcItemDynamic(s.ReadLineAutoCompleteResults()))),
				readline.PcItem("results",
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets(),
						readline.PcItemDynamic(s.ReadLineAutoCompleteModuleResults())))),
		),
		readline.PcItemDynamic(s.ReadLineAutoCompleteListModules(),
			readline.PcItem("list"),
			readline.PcItem("target",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("filter",
				readline.PcItemDynamic(s.ReadLineAutoCompleteFilters())),
			readline.PcItem("set",
				readline.PcItem("TARGET")),
			readline.PcItem("run"),
		),
		readline.PcItem("help"),
		readline.PcItem("info",
			readline.PcItem("session"),
			readline.PcItem("api")),
		readline.PcItem("api",
			readline.PcItem("run"),
			readline.PcItem("stop"),),
	)
	s.Prompt = &readline.Config{
		Prompt:          "\033[90m[OPF v"+s.Version+"]:\033[0m ",
		HistoryFile:     s.Config.Common.HistoryFile,
		InterruptPrompt: "^C",
		AutoComplete: completer,
		EOFPrompt:       "exit",
		HistorySearchFold:   true,
	}
}

func (s *Session) ParseCommand(line string){
	moduleName := strings.Split(line, " ")[0]
	module, err := s.SearchModule(moduleName)
	if err != nil{
		if moduleName == "help"{
			module, err = s.SearchModule("session_help")
		} else if !strings.HasPrefix(strings.TrimSpace(line), "target ") {
			s.Stream.Error("command '"+line+"' do not exist")
			s.Stream.Error("'help' for more information")
			return
		}
	}
	if strings.Contains(line, " "){
		if strings.HasPrefix(line, "target "){
			arguments := strings.Split(strings.TrimSpace(line), " ")
			switch arguments[1] {
			case "add":
				value := strings.SplitN(strings.TrimSpace(line), " ", 4)
				if len(arguments) < 4{
					s.Stream.Error("Please use subject add <type> <name>")
					return
				}
				id, err := s.AddTarget(value[2], value[3])
				if err != nil{
					s.Stream.Error(err.Error())
					return
				}
				s.Stream.Success("target '" + value[3] + "' as successfully added with id '"+id+"'")
			case "list":
				s.ListTargets()
			case "links":
				value := strings.SplitN(strings.TrimSpace(line), " ", 3)
				if len(arguments) < 3{
					s.Stream.Error("Please use subject add <type> <name>")
					return
				}
				trg, err := s.GetTarget(value[2])
				if err != nil{
					s.Stream.Error(err.Error())
					return
				}
				trg.Linked()
			case "view":
				if len(arguments) < 5{
					s.Stream.Error("Please use target view result <target_id> <result_id>")
					return
				}
				switch arguments[2] {
					case "results":
						value := strings.SplitN(strings.TrimSpace(line), " ", 5)
						if len(arguments) < 4{
							s.Stream.Error("Please use target view results <target_id>")
							return
						}
						trg, err := s.GetTarget(value[3])
						if err != nil{
							s.Stream.Error(err.Error())
							return
						}
						moduleName := value[4]
						results, err := trg.GetModuleResults(moduleName)
						if err != nil{
							s.Stream.Error(err.Error())
							return
						}

						t := s.Stream.GenerateTable()
						t.SetOutputMirror(os.Stdout)
						t.SetAllowedColumnLengths([]int{0, 30, 30, 30})
						headerRow := table.Row{}
						for _, result := range results{
							resRow := table.Row{}
							separator := trg.GetSeparator()
							header := strings.Split(result.Header, separator)
							res := strings.Split(result.Value, separator)
							if len(headerRow) < 1 {
								for _, h := range header {
									headerRow = append(headerRow, h)
								}
								headerRow = append(headerRow, "result_id")
								t.AppendHeader(headerRow)
							}
							for _, r := range res{
								resRow = append(resRow, r)
							}
							resRow = append(resRow, result.ResultId)
							t.AppendRow(resRow)
						}
						s.Stream.Render(t)

					case "result":
						value := strings.SplitN(strings.TrimSpace(line), " ", 5)
						trg, err := s.GetTarget(value[3])
						if err != nil{
							s.Stream.Error(err.Error())
							return
						}
						resultId := value[4]
						result, err := trg.GetResult(resultId)
						if err != nil{
							s.Stream.Error(err.Error())
							return
						}
						t := s.Stream.GenerateTable()
						t.SetOutputMirror(os.Stdout)
						separator := trg.GetSeparator()
						header := strings.Split(result.Header, separator)
						res := strings.Split(result.Value, separator)
						headerRow := table.Row{}
						resRow := table.Row{}
						for _, h := range header{
							headerRow = append(headerRow, h)
						}
						headerRow = append(headerRow, "RESULT ID")
						for _, r := range res{
							resRow = append(resRow, r)
						}
						resRow = append(resRow, result.ResultId)
						t.AppendHeader(headerRow)
						t.AppendRow(resRow)
						s.Stream.Render(t)
						return
				}

			case "update":
				value := strings.SplitN(strings.TrimSpace(line), " ", 4)
				if len(arguments) < 3{
					s.Stream.Error("Please use target update <target_id> <name>")
					return
				}
				s.UpdateTarget(value[2], value[3])
				s.Stream.Success("target '" + value[2] + "' as successfully updated.")
			case "modules":
				value := strings.SplitN(strings.TrimSpace(line), " ", 3)
				if len(arguments) < 3{
					s.Stream.Error("Please use target update <target_id> <name>")
					return
				}
				trg, err := s.GetTarget(value[2])
				if err != nil{
					s.Stream.Error(err.Error())
					return
				}

				t := s.Stream.GenerateTable()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{
					"Module",
					"Description",
					"Author",
					"Type",
				})
				for _, mod := range s.Modules{
					if mod.GetType() == trg.GetType() {
						t.AppendRow(table.Row{
							mod.Name(),
							mod.Description(),
							mod.Author(),
							mod.GetType(),
						})
					}
				}
				s.Stream.Render(t)
			case "delete":
				value := strings.SplitN(strings.TrimSpace(line), " ", 3)
				if len(arguments) < 3{
					s.Stream.Error("Please use target add <type> <name>")
					return
				}
				_, err := s.RemoveTarget(value[2])
				if err != nil{
					s.Stream.Error(err.Error())
					return
				}
				s.Stream.Success("target '" + value[2] + "' as successfully deleted.")

			}
		} else {
			arguments := strings.Split(strings.TrimSpace(line), " ")
			switch arguments[1] {
			case "target":
				if len(arguments) < 3 {
					s.Stream.Error("Please use <module> <set> <argument> <value>")
					return
				}
				ret, err := module.SetParameter("TARGET", arguments[2])
				if ret == false {
					s.Stream.Error(err.Error())
					return
				}
			case "filter":
				if len(arguments) < 3 {
					s.Stream.Error("Please use <module> <set> <argument> <value>")
					return
				}
				filter, errFilter := s.SearchFilter(arguments[2])
				if errFilter != nil{
					s.Stream.Error(errFilter.Error())
					return
				}
				if filter.WorkWith(arguments[0]) {
					ret, err := module.SetParameter("FILTER", arguments[2])
					if ret == false {
						s.Stream.Error(err.Error())
						return
					}
				} else{
					s.Stream.Error("This filter do not work with module '" + arguments[0] + "'")
					return
				}
			case "set":
				if len(arguments) < 4 {
					s.Stream.Error("Please use <module> <set> <argument> <value>")
					return
				}
				ret, err := module.SetParameter(arguments[2], arguments[3])
				if ret == false {
					s.Stream.Error(err.Error())
					return
				}
			case "list":
				module.ListArguments()
			case "run":
				if module.CheckRequired() {
					s.Information.ModuleLaunched = s.Information.ModuleLaunched + 1
					module.Start()
					filter, err := module.GetParameter("FILTER")
					if err == nil && filter.Value != ""{
						flt, err := s.SearchFilter(filter.Value)
						if err != nil{
							s.Stream.Error("Filter '"+filter.Value+"' as not found.")
							return
						}
						s.Stream.Success("Start filter '"+filter.Value+"'...")
						flt.Start(module)
					}
				} else {
					s.Stream.Error("Please validate required argument. (<module> list)")
				}
			}
		}
	}
	if moduleName == "help"{
		module.Start()
		filter, err := module.GetParameter("FILTER")
		if err == nil && filter.Value != ""{
			flt, err := s.SearchFilter(filter.Value)
			if err != nil{
				s.Stream.Error("Filter '"+filter.Value+"' as not found.")
				return
			}
			s.Stream.Success("Start filter '"+filter.Value+"'...")
			flt.Start(module)
		}
	}
}

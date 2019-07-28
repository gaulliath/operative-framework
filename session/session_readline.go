package session

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/graniet/go-pretty/table"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (s *Session) ClearScreen(){
	switch runtime.GOOS {
	case "linux", "darwin", "freebsd", "dragonfly", "netbsd", "openbsd", "solaris":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		// do nothing
	}
}

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
			if value[1] == "add"{
				for _, target := range s.Targets{
					for _, module := range target.Results{
						for _, result := range module{
							if result.Value != "" {
								returnResult = append(returnResult, result.ResultId)
							}
						}
					}
				}
				return returnResult
			}
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
			readline.PcItem("link",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets(),
					readline.PcItemDynamic(s.ReadLineAutoCompleteTargets()))),
			readline.PcItem("links",
				readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
			readline.PcItem("modules",readline.PcItemDynamic(s.ReadLineAutoCompleteTargets())),
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
		readline.PcItem("env"),
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

func (s *Session) ParseCommand(line string) []string{
	moduleName := strings.Split(line, " ")[0]
	module, err := s.SearchModule(moduleName)
	if err != nil{
		if moduleName == "help"{
			module, err = s.SearchModule("session_help")
		} else if !strings.HasPrefix(strings.TrimSpace(line), "target ") && !strings.HasPrefix(strings.TrimSpace(line), "note ") {
			s.Stream.Error("command '"+line+"' do not exist")
			s.Stream.Error("'help' for more information")
			return nil
		}
	}
	if strings.Contains(line, " "){
		if strings.HasPrefix(line, "sh "){
			arguments := strings.Split(strings.TrimSpace(line), " ")
			value := strings.SplitN(strings.TrimSpace(line), " ", 2)
			if len(arguments) < 2{
				s.Stream.Error("Please use sh <cmd> e.g: sh ls")
				return nil
			}
			_,_ = module.SetParameter("CMD", value[1])
			module.Start()
			return nil
		} else if strings.HasPrefix(line, "note "){
			arguments := strings.Split(strings.TrimSpace(line), " ")
			switch arguments[1] {
			case "add":
				value := strings.SplitN(strings.TrimSpace(line), " ", 4)
				findTarget, err := s.GetTarget(value[2])
				if err != nil{
					findResult, err := s.GetResult(value[2])
					if err != nil{
						s.Stream.Error("can't be find target/result with id '"+value[2]+"'")
						return nil
					}
					findResult.AddNote(value[3])
					s.Stream.Success("Note as been added to '"+value[2]+"'")
					return nil
				}
				findTarget.AddNote(value[3])
				s.Stream.Success("Note as been added to '"+value[2]+"'")
				return nil
			case "view":
				value := strings.SplitN(strings.TrimSpace(line), " ", 3)
				findTarget, err := s.GetTarget(value[2])
				t := s.Stream.GenerateTable()
				t.SetOutputMirror(os.Stdin)
				if err != nil{
					findResult, err := s.GetResult(value[2])
					fmt.Println(findResult)
					if err != nil{
						s.Stream.Error("can't be find target/result with id '"+value[2]+"'")
						return nil
					}
					t.AppendHeader(table.Row{
						"ID",
						"NOTE",
					})
					for _, note := range findResult.Notes{
						t.AppendRow(table.Row{
							note.Id,
							note.Text,
						})
					}
					s.Stream.Render(t)
					return nil
				}
				t.AppendHeader(table.Row{
					"ID",
					"NOTE",
				})
				for _, note := range findTarget.Notes{
					t.AppendRow(table.Row{
						note.Id,
						note.Text,
					})
				}
				s.Stream.Render(t)
				return nil
			}
		}else if strings.HasPrefix(line, "target "){
			arguments := strings.Split(strings.TrimSpace(line), " ")
			switch arguments[1] {
			case "add":
				value := strings.SplitN(strings.TrimSpace(line), " ", 4)
				if len(arguments) < 4{
					s.Stream.Error("Please use subject add <type> <name>")
					return nil
				}
				id, err := s.AddTarget(value[2], value[3])
				if err != nil{
					s.Stream.Error(err.Error())
					return nil
				}
				s.Stream.Success("target '" + value[3] + "' as successfully added with id '"+id+"'")
				return []string{
					id,
				}
			case "list":
				s.ListTargets()
			case "link":
				value := strings.SplitN(strings.TrimSpace(line), " ", 4)
				if len(arguments) < 3{
					s.Stream.Error("Please use subject add <type> <name>")
					return nil
				}
				trg, err := s.GetTarget(value[2])
				if err != nil{
					s.Stream.Error(err.Error())
					return nil
				}
				trg2, err := s.GetTarget(value[3])
				if err != nil{
					s.Stream.Error(err.Error())
					return nil
				}
				trg.Link(Linking{
					TargetId: trg2.GetId(),
				})
				s.Stream.Success("target '"+trg.GetId()+"' as linked to '"+trg2.GetId()+"'")
				return nil
			case "links":
				value := strings.SplitN(strings.TrimSpace(line), " ", 3)
				if len(arguments) < 3{
					s.Stream.Error("Please use subject links <target>")
					return nil
				}
				trg, err := s.GetTarget(value[2])
				if err != nil{
					s.Stream.Error(err.Error())
					return nil
				}
				trg.Linked()
			case "tag":
				switch arguments[2] {
					case "add":
						if len(arguments) < 5{
							s.Stream.Error("Please use target tag add <target_id> <tag>")
							return nil
						}
						value := strings.SplitN(strings.TrimSpace(line), " ", 5)
						trg, err := s.GetTarget(value[3])
						if err != nil{
							s.Stream.Error(err.Error())
							return nil
						}

						_, err = trg.AddTag(value[4])
						if err != nil {
							s.Stream.Error(err.Error())
							return nil
						}

						s.Stream.Success("Tag '"+value[4]+"' as been add to target '"+trg.GetName()+"'")
						return nil
					case "list":
						if len(arguments) < 4{
							s.Stream.Error("Please use target tag add <target_id> <tag>")
							return nil
						}
						value := strings.SplitN(strings.TrimSpace(line), " ", 4)
						trg, err := s.GetTarget(value[3])
						if err != nil{
							s.Stream.Error(err.Error())
							return nil
						}
						t := s.Stream.GenerateTable()
						t.SetOutputMirror(os.Stdout)
						t.SetAllowedColumnLengths([]int{40, 30, 30, 30})
						headerRow := table.Row{}
						resRow := table.Row{}
						headerRow = append(headerRow, "TAG")
						for _, tag := range trg.GetTags(){
							resRow = append(resRow, tag)
						}
						t.AppendHeader(headerRow)
						t.AppendRow(resRow)
						s.Stream.Render(t)
						return nil

				}
			case "view":
				if len(arguments) < 5{
					s.Stream.Error("Please use target view result <target_id> <result_id>")
					return nil
				}
				switch arguments[2] {
					case "results":
						value := strings.SplitN(strings.TrimSpace(line), " ", 5)
						if len(arguments) < 4{
							s.Stream.Error("Please use target view results <target_id>")
							return nil
						}
						trg, err := s.GetTarget(value[3])
						if err != nil{
							s.Stream.Error(err.Error())
							return nil
						}
						moduleName := value[4]
						results, err := trg.GetModuleResults(moduleName)
						if err != nil{
							s.Stream.Error(err.Error())
							return nil
						}

						t := s.Stream.GenerateTable()
						t.SetOutputMirror(os.Stdout)
						t.SetAllowedColumnLengths([]int{40, 30, 30, 30})
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
							return nil
						}
						resultId := value[4]
						result, err := trg.GetResult(resultId)
						if err != nil{
							s.Stream.Error(err.Error())
							return nil
						}
						t := s.Stream.GenerateTable()
						t.SetOutputMirror(os.Stdout)
						t.SetAllowedColumnLengths([]int{40, 30, 30, 30})
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
						return nil
				}

			case "update":
				value := strings.SplitN(strings.TrimSpace(line), " ", 4)
				if len(arguments) < 3{
					s.Stream.Error("Please use target update <target_id> <name>")
					return nil
				}
				s.UpdateTarget(value[2], value[3])
				s.Stream.Success("target '" + value[2] + "' as successfully updated.")
			case "modules":
				value := strings.SplitN(strings.TrimSpace(line), " ", 3)
				if len(arguments) < 3{
					s.Stream.Error("Please use target update <target_id> <name>")
					return nil
				}
				trg, err := s.GetTarget(value[2])
				if err != nil{
					s.Stream.Error(err.Error())
					return nil
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
					return nil
				}
				_, err := s.RemoveTarget(value[2])
				if err != nil{
					s.Stream.Error(err.Error())
					return nil
				}
				s.Stream.Success("target '" + value[2] + "' as successfully deleted.")

			}
		} else {
			arguments := strings.Split(strings.TrimSpace(line), " ")
			switch arguments[1] {
			case "target":
				if len(arguments) < 3 {
					s.Stream.Error("Please use <module> <set> <argument> <value>")
					return nil
				}
				ret, err := module.SetParameter("TARGET", arguments[2])
				if ret == false {
					s.Stream.Error(err.Error())
					return nil
				}
			case "filter":
				if len(arguments) < 3 {
					s.Stream.Error("Please use <module> <set> <argument> <value>")
					return nil
				}
				filter, errFilter := s.SearchFilter(arguments[2])
				if errFilter != nil{
					s.Stream.Error(errFilter.Error())
					return nil
				}
				if filter.WorkWith(arguments[0]) {
					ret, err := module.SetParameter("FILTER", arguments[2])
					if ret == false {
						s.Stream.Error(err.Error())
						return nil
					}
				} else{
					s.Stream.Error("This filter do not work with module '" + arguments[0] + "'")
					return nil
				}
			case "set":
				if len(arguments) < 4 {
					s.Stream.Error("Please use <module> <set> <argument> <value>")
					return nil
				}
				expl := strings.SplitN(line, " ", 4)
				ret, err := module.SetParameter(expl[2], expl[3])
				if ret == false {
					s.Stream.Error(err.Error())
					return nil
				}
			case "list":
				module.ListArguments()
			case "run":
				if module.CheckRequired() {
					if len(module.GetExternal()) > 0{
						for _, external := range module.GetExternal(){
							_, err := exec.LookPath(external)
							if err != nil {
								s.Stream.Error("This module need external program : '" + external + "'")
								return nil
							}
						}
					}
					s.Information.ModuleLaunched = s.Information.ModuleLaunched + 1
					background, errBack := module.GetParameter("BACKGROUND")
					if errBack == nil && strings.ToLower(background.Value) == "true"{
						go func(s *Session,m Module){
							s.Stream.Success("Running '" + module.Name() + "' in background...")
							module.Start()
							filter, err := module.GetParameter("FILTER")
							if err == nil && filter.Value != "" {
								flt, err := s.SearchFilter(filter.Value)
								if err != nil {
									s.Stream.Error("Filter '" + filter.Value + "' as not found.")
									return
								}
								s.Stream.Success("Start filter '" + filter.Value + "'...")
								flt.Start(module)
							}
							s.Stream.Success("Module '" + module.Name() + "' executed")
						}(s, module)
					} else {
						module.Start()
						r := module.GetResults()
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
						return r
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
				return nil
			}
			s.Stream.Success("Start filter '"+filter.Value+"'...")
			flt.Start(module)
		}
	}
	return nil
}

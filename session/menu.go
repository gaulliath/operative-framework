package session

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/torindex/helper"
	"os"
	"os/exec"
	"strings"
	"time"
)

func LoadAliasMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	switch arguments[1] {
	case "add":
		value := strings.SplitN(strings.TrimSpace(line), " ", 4)
		if len(value) < 4 {
			s.Stream.Error("Please use alias add <module> <alias> e.g: alias add google.twitter gt")
			return nil
		}
		module := value[2]
		s.AddAlias(value[3], module)
		return nil
	case "list":
		s.ListAlias()
		return nil
	}
	return nil
}

// @todo solve error in parsing line
func LoadNoteMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	switch strings.ToLower(arguments[1]) {
	case "add":
		if len(arguments) < 4 {
			s.Stream.Error("required argument are missing. E.g: note add <targetId/resultId> <text>")
			return nil
		}
		return nil
	case "delete":
		return nil
	case "view":
		return nil
	}
	return nil
}

func LoadTargetMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	switch arguments[1] {
	case "add":
		value := strings.SplitN(strings.TrimSpace(line), " ", 4)
		if len(arguments) < 4 {
			s.Stream.Error("Please use subject add <type> <name>")
			return nil
		}
		id, err := s.AddTarget(value[2], value[3])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}
		s.Stream.Success("target '" + value[3] + "' as successfully added with id '" + id + "'")
		return []string{
			id,
		}
	case "list":
		s.ListTargets()
	case "link":
		value := strings.SplitN(strings.TrimSpace(line), " ", 4)
		if len(arguments) < 4 {
			s.Stream.Error("Please use 'target link <target1> <target2>'")
			return nil
		}
		trg, err := s.GetTarget(value[2])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}
		trg2, err := s.GetTarget(value[3])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}
		if trg2.GetId() == trg.GetId() {
			s.Stream.Error("Can't link same target '" + trg.GetId() + "' : '" + trg2.GetId() + "'")
			return nil
		}
		trg.Link(Linking{
			TargetId: trg2.GetId(),
		})
		s.Stream.Success("target '" + trg.GetId() + "' as linked to '" + trg2.GetId() + "'")
		return nil
	case "links":
		value := strings.SplitN(strings.TrimSpace(line), " ", 3)
		if len(arguments) < 3 {
			s.Stream.Error("Please use subject links <target>")
			return nil
		}
		trg, err := s.GetTarget(value[2])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}
		trg.Linked()
	case "tag":
		switch arguments[2] {
		case "add":
			if len(arguments) < 5 {
				s.Stream.Error("Please use target tag add <target_id> <tag>")
				return nil
			}
			value := strings.SplitN(strings.TrimSpace(line), " ", 5)
			trg, err := s.GetTarget(value[3])
			if err != nil {
				s.Stream.Error(err.Error())
				return nil
			}

			_, err = s.AddTag(trg, value[4])
			if err != nil {
				s.Stream.Error(err.Error())
				return nil
			}

			s.Stream.Success("Tag '" + value[4] + "' as been add to target '" + trg.GetName() + "'")
			return nil
		case "list":
			if len(arguments) < 4 {
				s.Stream.Error("Please use target tag add <target_id> <tag>")
				return nil
			}
			value := strings.SplitN(strings.TrimSpace(line), " ", 4)
			trg, err := s.GetTarget(value[3])
			if err != nil {
				s.Stream.Error(err.Error())
				return nil
			}
			t := s.Stream.GenerateTable()
			t.SetOutputMirror(os.Stdout)
			t.SetAllowedColumnLengths([]int{40, 30, 30, 30})
			t.AppendHeader(table.Row{
				"TAG ID",
				"TEXT",
			})
			for _, tag := range trg.GetTags() {
				t.AppendRow(table.Row{
					tag.TagId,
					tag.Text,
				})
			}
			s.Stream.Render(t)
			return nil

		}
	case "view":
		if len(arguments) < 5 {
			s.Stream.Error("Please use target view result <target_id> <result_id>")
			return nil
		}
		switch arguments[2] {
		case "results":
			value := strings.SplitN(strings.TrimSpace(line), " ", 5)
			if len(arguments) < 4 {
				s.Stream.Error("Please use target view results <target_id>")
				return nil
			}
			trg, err := s.GetTarget(value[3])
			if err != nil {
				s.Stream.Error(err.Error())
				return nil
			}
			moduleName := value[4]
			results, err := trg.GetModuleResults(moduleName)
			if err != nil {
				s.Stream.Error(err.Error())
				return nil
			}

			t := s.Stream.GenerateTable()
			t.SetOutputMirror(os.Stdout)
			t.SetAllowedColumnLengths([]int{40, 30, 30, 30})
			headerRow := table.Row{}
			for _, result := range results {
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
				for _, r := range res {
					resRow = append(resRow, r)
				}
				resRow = append(resRow, result.ResultId)
				t.AppendRow(resRow)
			}
			s.Stream.Render(t)

		case "result":
			value := strings.SplitN(strings.TrimSpace(line), " ", 5)
			trg, err := s.GetTarget(value[3])
			if err != nil {
				s.Stream.Error(err.Error())
				return nil
			}
			resultId := value[4]
			result, err := trg.GetResult(resultId)
			if err != nil {
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
			for _, h := range header {
				headerRow = append(headerRow, h)
			}
			headerRow = append(headerRow, "RESULT ID")
			for _, r := range res {
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
		if len(arguments) < 3 {
			s.Stream.Error("Please use target update <target_id> <name>")
			return nil
		}
		s.UpdateTarget(value[2], value[3])
		s.Stream.Success("target '" + value[2] + "' as successfully updated.")

	case "modules":
		value := strings.SplitN(strings.TrimSpace(line), " ", 3)
		if len(arguments) < 3 {
			s.Stream.Error("Please use target update <target_id> <name>")
			return nil
		}
		trg, err := s.GetTarget(value[2])
		if err != nil {
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
		for _, mod := range s.Modules {
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
		if len(arguments) < 3 {
			s.Stream.Error("Please use target add <type> <name>")
			return nil
		}
		_, err := s.RemoveTarget(value[2])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}
		s.Stream.Success("target '" + value[2] + "' as successfully deleted.")
	}
	return nil
}

func LoadFindCommandMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	if len(arguments) < 3 {
		s.Stream.Error("Please use find / regex <search> <source> e.g: find operative framework results")
		return nil
	}

	search := strings.SplitN(strings.TrimSpace(line), " ", 2)[1]
	searchIn := arguments[len(arguments)-1]
	search = strings.TrimSpace(strings.Replace(search, searchIn, "", -1))

	_, _ = module.SetParameter("search", search)
	_, _ = module.SetParameter("source", searchIn)
	module.Start()
	return nil
}

func LoadShCommandMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	value := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(arguments) < 2 {
		s.Stream.Error("Please use sh <cmd> e.g: sh ls")
		return nil
	}
	_, _ = module.SetParameter("CMD", value[1])
	module.Start()
	return nil
}

func LoadModuleMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	switch strings.ToLower(arguments[1]) {
	case "target":
		if len(arguments) < 3 {
			s.Stream.Error("Please use <module> <set> <argument> <value>")
			return nil
		}

		_, err := s.GetTarget(arguments[2])
		if err != nil {
			newTarget := strings.SplitN(line, " ", 3)
			arguments[2], err = s.AddTarget(module.GetType(), newTarget[2])
			if err != nil {
				return nil
			}
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
		if errFilter != nil {
			s.Stream.Error(errFilter.Error())
			return nil
		}
		if filter.WorkWith(arguments[0]) {
			ret, err := module.SetParameter("FILTER", arguments[2])
			if ret == false {
				s.Stream.Error(err.Error())
				return nil
			}
		} else {
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
		return nil
	case "reset:target":
		_, _ = module.SetParameter("TARGET", "")
		return nil
	case "run":
		if module.CheckRequired() {
			if len(module.GetExternal()) > 0 {
				for _, external := range module.GetExternal() {
					_, err := exec.LookPath(external)
					if err != nil {
						s.Stream.Error("This module need external program : '" + external + "'")
						return nil
					}
				}
			}
			s.Information.ModuleLaunched = s.Information.ModuleLaunched + 1
			background, errBack := module.GetParameter("BACKGROUND")
			if errBack == nil && strings.ToLower(background.Value) == "true" {
				go func(s *Session, m Module) {
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
				startedAt := time.Now()
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

				if s.Config.PushDriver == "ONLY_SERVER" ||
					s.Config.PushDriver == "SCREEN_AND_SERVER" {

					targetId, err := module.GetParameter("TARGET")
					if err != nil {
						s.Stream.Error("error with a push to webserver: " + err.Error())
						return nil
					}

					target, err := s.GetTarget(targetId.Value)
					if err != nil {
						s.Stream.Error("error with a push to webserver: " + err.Error())
						return nil
					}

					if _, ok := target.Results[module.Name()]; ok {
						s.PushResultsToGate(target.Results[module.Name()], startedAt)
					}
				}
				return r
			}
		} else {
			s.Stream.Error("Please validate required argument. (" + module.Name() + " list)")
		}
	}
	return nil
}

func LoadIntervalCommandMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	switch arguments[1] {
	case "generate":
		if len(arguments) < 3 {
			s.Stream.Error("Please use interval generate command1;command2...")
			return nil
		}
		command := strings.SplitN(line, " ", 3)
		newInterval := s.GenerateInterval(command[2])
		s.Stream.Success("new interval as generated with id '" + newInterval.Id + "'")
		break
	case "list":
		t := s.Stream.GenerateTable()
		t.SetOutputMirror(os.Stdout)
		t.SetAllowedColumnLengths([]int{30, 30, 30})
		t.AppendHeader(table.Row{
			"ID",
			"ACTIVATE",
			"COMMAND",
			"DELAY",
			"LAST",
			"NEXT",
		})

		for _, interval := range s.Interval {
			ActivatedString := color.RedString("false")
			if interval.Activated {
				ActivatedString = color.GreenString("true")
			}
			t.AppendRow(table.Row{
				interval.Id,
				ActivatedString,
				interval.Commands,
				interval.Delay,
				interval.LastRun.Format("2006-01-02 15:04:05"),
				interval.NextRun.Format("2006-01-02 15:04:05"),
			})
		}

		s.Stream.Render(t)
		break
	case "set":
		if len(arguments) < 5 {
			s.Stream.Error("Please use : interval set <intervalId> <argument> <value>")
			return nil
		}
		commands := strings.SplitN(line, " ", 5)
		interval, err := s.GetInterval(commands[2])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}
		switch strings.ToLower(commands[3]) {
		case "delay":
			interval.SetDelay(helper.StringToInt(commands[4]))
			break
		case "command":
			interval.SetCommand(commands[4])
			break
		}
		s.Stream.Success("command executed.")
		break
	case "up":
		if len(arguments) < 3 {
			s.Stream.Error("Please use : interval run <intervalId>")
			return nil
		}
		options := strings.SplitN(line, " ", 3)
		interval, err := s.GetInterval(options[2])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}

		if interval.GetDelay() == 0 {
			s.Stream.Error("Please configure a delay. E.g: interval set <intervalId> DELAY 10")
			return nil
		}
		interval.Up()
		s.Stream.Success("interval '" + interval.Id + "' as started at '" + time.Now().Format("2006-01-02 15:04:05") + "'")
		break
	case "down":
		if len(arguments) < 3 {
			s.Stream.Error("Please use : interval run <intervalId>")
			return nil
		}
		options := strings.SplitN(line, " ", 3)
		interval, err := s.GetInterval(options[2])
		if err != nil {
			s.Stream.Error(err.Error())
			return nil
		}
		interval.Down()
		s.Stream.Success("interval '" + interval.Id + "' as stopped at '" + time.Now().Format("2006-01-02 15:04:05") + "'")
		break
	}

	return nil
}

func LoadModuleByTypeMenu(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	if len(arguments) < 2 {
		s.ListModules()
		return nil
	}

	stype := arguments[1]
	if !s.CheckTypeExist(stype) {
		s.Stream.Error("Please select a valid type : " + strings.Join(s.ListType(), ","))
		return nil
	}
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"NAME",
		"DESCRIPTION",
	})
	for _, module := range s.Modules {
		if strings.ToLower(module.GetType()) == strings.ToLower(stype) {
			t.AppendRow(table.Row{
				module.Name(),
				module.Description(),
			})
		}
	}
	s.Stream.Render(t)
	return nil
}

func LoadAnalyticsWebBased(line string, module Module, s *Session) []string {
	arguments := strings.Split(strings.TrimSpace(line), " ")
	switch strings.ToLower(arguments[1]) {
	case "up":
		a, al := s.GenerateAnalytics()
		aJson, _ := json.Marshal(a)
		alJson, _ := json.Marshal(al)

		fmt.Println(string(aJson))
		fmt.Println("====")
		fmt.Println(string(alJson))
		break
	}

	return nil
}

func LoadEventsMenu(line string, module Module, s *Session) []string {
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{0, 30})
	t.AppendHeader(table.Row{
		"ID",
		"TYPE",
		"VALUE",
	})

	for _, event := range s.Events {
		t.AppendRow(table.Row{
			event.EventId,
			event.Type,
			event.Value,
		})
	}
	s.Stream.Render(t)
	return nil
}

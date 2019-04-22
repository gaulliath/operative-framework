package session

import (
	"github.com/chzyer/readline"
	"github.com/graniet/go-pretty/table"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type Session struct{
	Id int `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionName string `json:"session_name"`
	Connection Connection `json:"-" sql:"-"`
	Version string            `json:"version" sql:"-"`
	Targets []*Target   `json:"subjects" sql:"-"`
	Modules []Module          `json:"modules" sql:"-"`
	Prompt *readline.Config `json:"-" sql:"-"`
	Stream Stream `json:"-" sql:"-"`
}

func (Session) TableName() string{
	return "sessions"
}


func (s *Session) GetId() int{
	return s.Id
}

func (s *Session) ListType() []string{
	return []string{
		"enterprise",
		"ip_address",
		"website",
		"url",
		"person",
		"social_network",
	}
}

func (s *Session) GetTarget(id string) (*Target, error){
	for _, targ := range s.Targets{
		if targ.GetId() == id{
			return targ, nil
		}
	}
	return nil, errors.New("can't find selected target")
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
					module.Start()
				} else {
					s.Stream.Error("Please validate required argument. (<module> list)")
				}
			}
		}
	}
	if moduleName == "help"{
		module.Start()
	}
}

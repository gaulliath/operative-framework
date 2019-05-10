package session

import (
	"errors"
	"github.com/graniet/go-pretty/table"
	"os"
	"strconv"
	"strings"
)

const (
	INT = 1
	STRING = 2
	BOOL = 3
	FLOAT = 4
)

type Module interface {
	Start()
	Name() string
	Author() string
	Description() string
	GetType() string
	ListArguments()
	GetExport() []TargetResults
	SetExport(result TargetResults)
	GetResults() []string
	GetInformation() ModuleInformation
	CheckRequired() bool
	SetParameter(name string, value string) (bool, error)
	GetParameter(name string) (Param, error)
	GetAllParameters() []Param
	WithProgram(name string) bool
	GetExternal() []string
	CreateNewParam(name string, description string, value string, isRequired bool, paramType int)
}

type Param struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Value string `json:"value"`
	IsRequired bool `json:"is_required"`
	ParamType int `json:"param_type"`
}

type SessionModule struct{
	Module
	Export []TargetResults
	Parameters []Param `json:"parameters"`
	History []string `json:"history"`
	External []string `json:"external"`
	Results []string
}

type ModuleInformation struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Author string `json:"author"`
	Type string `json:"type"`
	Parameters []Param `json:"parameters"`
}

type TargetResults struct{
	Id int `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionId int `json:"-" gorm:"session_id"`
	ModuleName string `json:"module_name"`
	ResultId string `json:"result_id" gorm:"primary_key:yes;column:result_id"`
	TargetId string `json:"target_id" gorm:"target_id"`
	Header string `json:"key" gorm:"result_header"`
	Value string `json:"value" gorm:"result_value"`
	Notes []Note
}

func (result *TargetResults) AddNote(text string){
	result.Notes = append(result.Notes, Note{
		Text: text,
	})
	return
}

func (s *Session) SearchModule(name string)(Module, error){
	for _, module := range s.Modules{
		if module.Name() == name{
			return module, nil
		}
	}
	return nil, errors.New("error: This module not found")
}


func (module *SessionModule) GetParameter(name string) (Param, error){
	for _, param := range module.Parameters{
		if param.Name == name{
			return param, nil
		}
	}
	return Param{}, errors.New("parameter not found")
}

func (module *SessionModule) SetParameter(name string, value string) (bool, error){
	for k, param := range module.Parameters{
		if param.Name == name{
			module.Parameters[k].Value = value
			return true, nil
		}
	}
	return false, errors.New("argument not found")
}

func (module *SessionModule) CheckRequired() bool{
	for _, param := range module.Parameters{
		if param.IsRequired == true{
			switch param.ParamType {
			case STRING:
				if param.Value == "" {
					return false
				}
			case INT:
				value, _ := strconv.Atoi(param.Value)
				if value == 0 {
					return false
				}
			case BOOL:
				value := strings.TrimSpace(param.Value)
				if value == ""{
					return false
				}
			}
		}
	}
	return true
}

func (module *SessionModule) CreateNewParam(name string, description string, value string, isRequired bool, paramType int){
	newParam := Param{
		Name:strings.ToUpper(name),
		Value:value,
		Description: description,
		IsRequired:isRequired,
		ParamType:paramType,
	}
	module.Parameters = append(module.Parameters, newParam)
}

func (module *SessionModule) WithProgram(name string) bool{
	module.External = append(module.External, name)
	return true
}

func (module *SessionModule) ListArguments(){
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"argument", "description" ,"value", "required", "type"})
	if len(module.Parameters) > 0{
		for _, argument := range module.Parameters{
			argumentType := ""
			argumentRequired := ""

			if argument.ParamType == STRING{
				argumentType = "STRING"
			} else if argument.ParamType == INT{
				argumentType = "INTEGER"
			} else if argument.ParamType == BOOL{
				argumentType = "BOOLEAN"
			}

			if argument.IsRequired == true{
				argumentRequired = "YES"
			} else{
				argumentRequired = "NO"
			}

			if argument.Value == ""{
				argument.Value = "NO DEFAULT"
			}
			t.AppendRow([]interface{}{argument.Name, argument.Description, argument.Value, argumentRequired, argumentType})
		}
	} else{
		t.AppendRow([]interface{}{"No argument."})
	}
	t.Render()
}


func (module *SessionModule) SetExport(result TargetResults){
	module.Export = append(module.Export, result)
}

func (module *SessionModule) GetExport() []TargetResults{
	return module.Export
}

func (module *SessionModule) GetAllParameters() []Param{
	return module.Parameters
}

func (module *SessionModule) GetResults() []string{
	return module.Results
}

func (module *SessionModule) GetExternal() []string{
	return module.External
}

package session

import (
	"encoding/base64"
	"github.com/graniet/go-pretty/table"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"os"
)

type Target struct{
	Id int `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionId int `json:"-" gorm:"column:session_id"`
	TargetId string `json:"id" gorm:"column:target_id"`
	Sess *Session `json:"-" gorm:"-"`
	Name string `json:"name" gorm:"column:target_name"`
	Type string `json:"type" gorm:"column:target_type"`
	Results map[string][]TargetResults `sql:"-"`
	TargetLinked []Linking `json:"target_linked" sql:"-"`
}

type Linking struct{
	LinkingId int `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	TargetBase string `json:"target_base"`
	TargetId string `json:"target_id"`
	TargetName string `json:"target_name"`
	TargetType string `json:"target_type"`
	TargetResultId string `json:"target_result_id"`
}

func (Linking) TableName() string{
	return "target_links"
}

func (sub *Target) GetId() string{
	return sub.TargetId
}

func (sub *Target) GetName() string{
	return sub.Name
}

func (sub *Target) GetType() string{
	return sub.Type
}

func (sub *Target) GetResults() map[string][]TargetResults{
	return sub.Results
}

func (sub *Target) GetLinked() []Linking{
	return sub.TargetLinked
}

func (sub *Target) PushLinked(t Linking){
	sub.TargetLinked = append(sub.TargetLinked, t)
}

func (target *Target) CheckType() bool{
	for _, sType := range target.Sess.ListType(){
		if sType == target.GetType(){
			return true
		}
	}
	return false
}

func (target *Target) Link(target2 Linking){
	if target.GetId() == target2.TargetId{
		return
	}
	t2, err := target.Sess.GetTarget(target2.TargetId)
	if err != nil{
		target.Sess.Stream.Error(err.Error())
		return
	}

	for _, trg := range target.TargetLinked{
		if trg.TargetId == t2.GetId(){
			return
		}
	}
	target2.TargetType = t2.GetType()
	target2.TargetName = t2.GetName()
	target2.TargetBase = target.GetId()
	target.PushLinked(target2)
	t2.PushLinked(Linking{
		TargetBase: t2.GetId(),
		TargetName: target.GetName(),
		TargetId: target.GetId(),
		TargetType: target.GetType(),
		TargetResultId: target2.TargetResultId,
	})
	target.Sess.Connection.ORM.Create(&target2)
	target.Sess.Connection.ORM.Create(&Linking{
		TargetBase: t2.GetId(),
		TargetName: target.GetName(),
		TargetId: target.GetId(),
		TargetType: target.GetType(),
		TargetResultId: target2.TargetResultId,
	})
}

func (target *Target) GetResult(id string) (TargetResults, error){
	for _, module := range target.Results{
		for _, result := range module{
			if result.ResultId == id{
				return result, nil
			}
		}
	}
	return TargetResults{}, errors.New("Result as been not found.")
}

func (target *Target) Linked(){
	t := target.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"TARGET",
		"NAME",
		"TYPE",
		"RESULT ID",

	})
	for _, element := range target.TargetLinked{
		t.AppendRow(table.Row{
			element.TargetId,
			element.TargetName,
			element.TargetType,
			element.TargetResultId,
		})
	}
	target.Sess.Stream.Render(t)
}

func (target *Target) GetSeparator() string{
	return base64.StdEncoding.EncodeToString([]byte(";operativeframewor;"))[0:5]
}

func (target *Target) Save(module Module, result TargetResults) bool{
	result.ResultId = ksuid.New().String()
	result.TargetId = target.GetId()
	result.SessionId = target.Sess.GetId()
	result.ModuleName = module.Name()
	target.Results[module.Name()] = append(target.Results[module.Name()], result)
	target.Sess.Connection.ORM.Create(&result).Table("target_results")
	targets, err := target.Sess.FindLinked(module.Name(), result)
	if err == nil {
		for _, id := range targets {
			target.Link(Linking{
				TargetId: id,
				TargetResultId: result.ResultId,
			})
		}
	}
	module.SetExport(result)
	return true
}
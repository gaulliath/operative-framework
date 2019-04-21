package session

import (
	"github.com/graniet/go-pretty/table"
	"os"
)

type Target struct{
	Id string `json:"id" gorm:"primary_key:yes;column:target_id"`
	Sess *Session `json:"-" gorm:"-"`
	Name string `json:"name" gorm:"column:target_name"`
	Type string `json:"type" gorm:"column:target_type"`
	Results map[string][]interface{} `sql:"-"`
	TargetLinked []*Target `json:"target_linked" sql:"-"`
}

func (sub *Target) GetId() string{
	return sub.Id
}

func (sub *Target) GetName() string{
	return sub.Name
}

func (sub *Target) GetType() string{
	return sub.Type
}

func (sub *Target) GetResults() map[string][]interface{}{
	return sub.Results
}

func (sub *Target) GetLinked() []*Target{
	return sub.TargetLinked
}

func (sub *Target) PushLinked(t *Target){
	sub.TargetLinked = append(sub.TargetLinked, t)
}

func (target *Target) ListType() []string{
	return []string{
		"enterprise",
		"ip_address",
		"website",
		"url",
		"person",
		"social_network",
	}
}

func (target *Target) CheckType() bool{
	for _, sType := range target.ListType(){
		if sType == target.GetType(){
			return true
		}
	}
	return false
}

func (target *Target) Link(target2 string){
	if target.GetId() == target2{
		return
	}
	t2, err := target.Sess.GetTarget(target2)
	if err != nil{
		target.Sess.Stream.Error(err.Error())
		return
	}

	for _, trg := range target.TargetLinked{
		if trg.GetId() == t2.GetId(){
			return
		}
	}
	target.PushLinked(t2)
	t2.PushLinked(target)
}

func (target *Target) Linked(){
	t := target.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"TARGET",
		"NAME",

	})
	for _, element := range target.TargetLinked{
		t.AppendRow(table.Row{
			element.GetId(),
			element.GetName(),
		})
	}
	target.Sess.Stream.Render(t)
}

func (target *Target) Save(module Module, result interface{}) bool{
	target.Results[module.Name()] = append(target.Results[module.Name()], result)
	targets, err := target.Sess.FindLinked(module.Name(), result)
	if err == nil {
		for _, id := range targets {
			target.Link(id)
		}
	}
	return true
}
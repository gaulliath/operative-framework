package session

import (
	"errors"
	"github.com/graniet/go-pretty/table"
	"github.com/segmentio/ksuid"
	"os"
	"strings"
)


func (s *Session) GetTarget(id string) (*Target, error){
	for _, targ := range s.Targets{
		if targ.GetId() == id{
			return targ, nil
		}
	}
	return nil, errors.New("can't find selected target")
}

func (s *Session) AddTarget(t string, name string) (string, error){
	subject := Target{
		SessionId: s.GetId(),
		TargetId: ksuid.New().String(),
		Name: name,
		Type: t,
		Results: make(map[string][]TargetResults),
		Sess: s,
	}
	if !subject.CheckType(){
		t := s.Stream.GenerateTable()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"TYPE",
		})
		for _, sType := range s.ListType(){
			t.AppendRow(table.Row{
				sType,
			})
		}
		s.Stream.Render(t)
		s.Stream.Error("Please configure valid type")
		return "", errors.New("this '"+subject.GetType()+"' type isn't available")
	}
	for _, subject := range s.Targets{
		if subject.GetName() == name && subject.GetType() == t{
			return "", errors.New("this target already exist on current session")
		}
	}
	s.Targets = append(s.Targets, &subject)
	s.Connection.ORM.Create(&subject)
	s.FindLinkedTargetByResult(&subject)
	return subject.GetId(), nil
}

func (s *Session) RemoveTarget(id string) (bool, error){
	if len(s.Targets) < 1{
		return false, errors.New("at this moment a session don't have target")
	}
	var newSubject []*Target
	for _, subject := range s.Targets{
		if subject.GetId() != id{
			var newLinked []Linking
			for _, linked := range subject.TargetLinked{
				if linked.TargetId != id{
					newLinked = append(newLinked, linked)
				}
			}
			subject.TargetLinked = newLinked
			newSubject = append(newSubject, subject)
		}
	}
	t, err := s.GetTarget(id)
	if err == nil{
		s.Connection.ORM.Delete(t)
	}
	s.Targets = newSubject
	return true, nil
}

func (s *Session) ListTargets(){
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "NAME", "TYPE", "MODULE RUN", "RESULTS"})
	for _, subject := range s.Targets{
		result := 0
		for k, _ := range subject.Results{
			result = result + len(subject.Results[k])
		}
		t.AppendRow(table.Row{
			subject.GetId(),
			subject.GetName(),
			subject.GetType(),
			len(subject.Results),
			result,
		})
	}
	s.Stream.Render(t)
}

func (s *Session) UpdateTarget(id string, value string){
	for k, t := range s.Targets{
		if t.GetId() == id{
			s.Targets[k].Name = value
			s.Connection.ORM.Save(t)
		}
	}
}

func (s *Session) FindLinked(m string, res TargetResults) ([]string, error){
	var targets []string
	for _,t := range s.Targets{
		targetId := t.GetId()
		for _, targetRes := range t.Results[m]{
			if res.Header == targetRes.Header && res.Value == targetRes.Value{
				targets = append(targets, targetId)
			}
		}
	}
	if len(targets) < 1 {
		return nil, errors.New("can't find linked target")
	} else{
		return targets, nil
	}
}

func (s *Session) FindLinkedTargetByResult(t *Target){
	targets := make(map[string]string, 0)
	for _, target := range s.Targets{
		for _, module := range target.Results{
			for _, res := range module{
				result := strings.Split(res.Value, target.GetSeparator())
				for _, r := range result{
					if strings.TrimSpace(strings.ToLower(r)) == strings.TrimSpace(strings.ToLower(t.Name)){
						targets[target.TargetId] = res.ResultId
					}
				}
			}
		}
	}
	if len(targets) > 0 {
		for TargetId, resId := range targets {
			trg, err := s.GetTarget(TargetId)
			if err == nil {
				trg.Link(Linking{
					TargetId:       t.TargetId,
					TargetResultId: resId,
				})
			}
		}
	}
}
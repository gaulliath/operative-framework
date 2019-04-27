package session

import (
	"errors"
)

type ModuleFilter interface {
	Start(mod Module)
	Name() string
	Description() string
	Author() string
	WorkWith(name string) bool
}

type SessionFilter struct{
	ModuleFilter
	With []string
}

func (s *Session) SearchFilter(name string)(ModuleFilter, error){
	for _, filter := range s.Filters{
		if filter.Name() == name{
			return filter, nil
		}
	}
	return nil, errors.New("error: This filter not found")
}

func (filter *SessionFilter) AddModule(name string){
	filter.With = append(filter.With, name)
	return
}

func (filter *SessionFilter) WorkWith(name string) bool{
	for _, module := range filter.With{
		if module == name{
			return true
		}
	}
	return false
}

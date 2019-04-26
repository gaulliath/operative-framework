package session

import (
	"errors"
)

type ModuleFilter interface {
	Start(mod Module)
	Name() string
	Description() string
	Author() string
}

type SessionFilter struct{
	ModuleFilter
}

func (s *Session) SearchFilter(name string)(ModuleFilter, error){
	for _, filter := range s.Filters{
		if filter.Name() == name{
			return filter, nil
		}
	}
	return nil, errors.New("error: This module not found")
}

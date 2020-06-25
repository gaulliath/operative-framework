package find

import (
	"github.com/graniet/operative-framework/session"
)

type FindModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushFindModule(s *session.Session) *FindModule {
	mod := FindModule{
		sess:   s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("search", "search term e.g: Operative Framework", "", true, session.STRING)
	mod.CreateNewParam("source", "or would you like to search?", "results", false, session.STRING)
	return &mod
}

func (module *FindModule) Name() string {
	return "find"
}

func (module *FindModule) Description() string {
	return "Find by search term in result(s) or target(s)"
}

func (module *FindModule) Author() string {
	return "Tristan Granier"
}

func (module *FindModule) GetType() []string {
	return []string{
		session.T_TARGET_COMMAND,
	}
}

func (module *FindModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *FindModule) Start() {
	return
}

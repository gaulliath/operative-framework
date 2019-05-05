package web_archive

import "github.com/graniet/operative-framework/session"

type WebArchiveModule struct{
	session.SessionModule
	sess *session.Session
	Stream session.Stream

}

func (module *WebArchiveModule) Name() string{
	return "web_archive"
}

func (module *WebArchiveModule) Author() string{
	return "Tristan Granier"
}

func (module *WebArchiveModule) Description() string{
	return "Search possible archive on site web"
}

func (module *WebArchiveModule) GetType() string{
	return "website"
}

func (module *WebArchiveModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

package modules

import (
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/operative-framework/modules/module_base/session_help"
	"github.com/graniet/operative-framework/modules/header_retrieval"
	"github.com/graniet/operative-framework/modules/linkedin_search"
	"github.com/graniet/operative-framework/modules/metatag_spider"
	"github.com/graniet/operative-framework/modules/module_base/session_stream"
	"github.com/graniet/operative-framework/modules/bing_vhost"
)

func LoadModules(s *session.Session){
	s.Modules = append(s.Modules, session_help.PushModuleHelp(s))
	s.Modules = append(s.Modules, header_retrieval.PushModuleHeaderRetrieval(s))
	s.Modules = append(s.Modules, linkedin_search.PushLinkedinSearchModule(s))
	s.Modules = append(s.Modules, metatag_spider.PushMetaTagModule(s))
	s.Modules = append(s.Modules, session_stream.PushSessionStreamModule(s))
	s.Modules = append(s.Modules, bing_vhost.PushBingVirtualHostModule(s))
}

package modules

import (
	"github.com/graniet/operative-framework/modules/bing_vhost"
	"github.com/graniet/operative-framework/modules/header_retrieval"
	"github.com/graniet/operative-framework/modules/instagram"
	"github.com/graniet/operative-framework/modules/linkedin_search"
	"github.com/graniet/operative-framework/modules/metatag_spider"
	"github.com/graniet/operative-framework/modules/module_base/session_help"
	"github.com/graniet/operative-framework/modules/module_base/session_stream"
	"github.com/graniet/operative-framework/modules/phone_generator"
	"github.com/graniet/operative-framework/modules/viewdns_search"
	"github.com/graniet/operative-framework/modules/whatsapp"
	"github.com/graniet/operative-framework/session"
)

func LoadModules(s *session.Session){
	s.Modules = append(s.Modules, session_help.PushModuleHelp(s))
	s.Modules = append(s.Modules, header_retrieval.PushModuleHeaderRetrieval(s))
	s.Modules = append(s.Modules, linkedin_search.PushLinkedinSearchModule(s))
	s.Modules = append(s.Modules, metatag_spider.PushMetaTagModule(s))
	s.Modules = append(s.Modules, session_stream.PushSessionStreamModule(s))
	s.Modules = append(s.Modules, bing_vhost.PushBingVirtualHostModule(s))
	s.Modules = append(s.Modules, viewdns_search.PushWSearchModule(s))
	s.Modules = append(s.Modules, phone_generator.PushPhoneGeneratorModule(s))
	s.Modules = append(s.Modules, whatsapp.PushWhatsappExtractorModule(s))
	s.Modules = append(s.Modules, instagram.PushInstagramFollowersModule(s))

	for _, mod := range s.Modules{
		mod.CreateNewParam("FILTER", "Use module filter after execution", "",false, session.STRING)
	}
}

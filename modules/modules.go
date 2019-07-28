package modules

import (
	"github.com/graniet/operative-framework/modules/account_checker"
	"github.com/graniet/operative-framework/modules/bing_vhost"
	"github.com/graniet/operative-framework/modules/darksearch"
	"github.com/graniet/operative-framework/modules/directory_search"
	"github.com/graniet/operative-framework/modules/get_ipaddress"
	"github.com/graniet/operative-framework/modules/google_dork"
	"github.com/graniet/operative-framework/modules/google_search"
	"github.com/graniet/operative-framework/modules/header_retrieval"
	"github.com/graniet/operative-framework/modules/image_reverse_search"
	"github.com/graniet/operative-framework/modules/instagram"
	"github.com/graniet/operative-framework/modules/linkedin_search"
	"github.com/graniet/operative-framework/modules/mac_vendor"
	"github.com/graniet/operative-framework/modules/metatag_spider"
	"github.com/graniet/operative-framework/modules/module_base/session_help"
	"github.com/graniet/operative-framework/modules/module_base/session_stream"
	"github.com/graniet/operative-framework/modules/pastebin"
	"github.com/graniet/operative-framework/modules/pastebin_email"
	"github.com/graniet/operative-framework/modules/phone_buster"
	"github.com/graniet/operative-framework/modules/phone_generator"
	"github.com/graniet/operative-framework/modules/phone_generator_fr"
	"github.com/graniet/operative-framework/modules/report_pdf"
	"github.com/graniet/operative-framework/modules/sample"
	"github.com/graniet/operative-framework/modules/searchsploit"
	"github.com/graniet/operative-framework/modules/societe_com"
	"github.com/graniet/operative-framework/modules/system"
	"github.com/graniet/operative-framework/modules/tools_suggester"
	"github.com/graniet/operative-framework/modules/twitter"
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
	s.Modules = append(s.Modules, instagram.PushInstagramFeedModule(s))
	s.Modules = append(s.Modules, instagram.PushInstagramFollowingModule(s))
	s.Modules = append(s.Modules, instagram.PushInstagramFriendsModule(s))
	s.Modules = append(s.Modules, instagram.PushInstagramInfoModule(s))
	s.Modules = append(s.Modules, twitter.PushTwitterFollowerModule(s))
	s.Modules = append(s.Modules, twitter.PushTwitterRetweetModule(s))
	s.Modules = append(s.Modules, twitter.PushTwitterFollowingModule(s))
	s.Modules = append(s.Modules, image_reverse_search.PushImageReverseModule(s))
	s.Modules = append(s.Modules, societe_com.PushSocieteComModuleModule(s))
	s.Modules = append(s.Modules, pastebin_email.PushPasteBinEmailModule(s))
	s.Modules = append(s.Modules, phone_buster.PushPhoneBusterModule(s))
	s.Modules = append(s.Modules, directory_search.PushModuleDirectorySearch(s))
	s.Modules = append(s.Modules, tools_suggester.PushModuleToolsSuggester(s))
	s.Modules = append(s.Modules, sample.PushSampleModuleModule(s))
	s.Modules = append(s.Modules, system.PushSystemModuleModule(s))
	s.Modules = append(s.Modules, pastebin.PushPasteBinModule(s))
	s.Modules = append(s.Modules, searchsploit.PushSearchSploitModule(s))
	s.Modules = append(s.Modules, get_ipaddress.PushGetIpAddressModule(s))
	s.Modules = append(s.Modules, report_pdf.PushReportPDFModule(s))
	s.Modules = append(s.Modules, account_checker.PushAccountCheckerModule(s))
	s.Modules = append(s.Modules, google_search.PushGoogleSearchModule(s))
	s.Modules = append(s.Modules, phone_generator_fr.PushPhoneGeneratorFrModule(s))
	s.Modules = append(s.Modules, mac_vendor.PushMacVendorModule(s))
	s.Modules = append(s.Modules, darksearch.PushMacVendorModule(s))
	s.Modules = append(s.Modules, google_dork.PushGoogleDorkModule(s))

	for _, mod := range s.Modules{
		s.PushType(mod.GetType())
		mod.CreateNewParam("FILTER", "Use module filter after execution", "",false, session.STRING)
		mod.CreateNewParam("BACKGROUND", "Run this task in background", "false", false, session.BOOL)
	}
}

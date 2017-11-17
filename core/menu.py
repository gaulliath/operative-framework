#!/usr/bin/env	python

import sys

menu_list = {	
	"quit":"exit_operative",
	"modules":"show_module",
	"browser_hack": "browser_hacks",
	"help":"show_help",
	"set":"set_enterprise",
	"run":"run_enterprise",
	"load_db":"load_db",
	"search_db":"search_dbs",
    "update":"update_framework",
	"clear":"clear_screen",
	"search_domain":"domain_module",
	"json_api": "load_api_json",
	"webserver": 'run_webserver',
	"campaign":"start_campaign",
	"new_module":"generate_module_class",
	"social_network":"social_network_gathering",
}

menu_shortcut = {
	"--campaign":"start_campaign",
	"--update":"update_framework",
	"--modules":"show_module",
	"--version":"banner",
    "--generate_module":"generate_module_class",
    "--use":"shortcut_load_module",
	"--api": "load_api_json",
	"--web" : "run_webserver",
	"--social_networl":"social_network_gathering",
}

menu_export = {
	"XML":"core/exports/XML",
	"JSON":"core/exports/JSON",
	"HTML":"core/exports/HTML"
}


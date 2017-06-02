#!/usr/bin/env	python

from colorama import Fore,Back,Style

import os, sys

def show_options(require):
    #print Back.WHITE + Fore.WHITE + "Module parameters" + Style.RESET_ALL
    for line in require:
        if require[line][0]["value"] == "":
	    value = "No value"
	else:
	    value = require[line][0]["value"]
	if require[line][0]["required"] == "yes":
	    if require[line][0]["value"] != "":
	        print Fore.GREEN+Style.BRIGHT+ "+ " +Style.RESET_ALL+line+ ": " +value
	    else:
	        print Fore.RED+Style.BRIGHT+ "- " +Style.RESET_ALL+line+ "(" +Fore.RED+ "is_required" +Style.RESET_ALL+ "):" +value
	else:
	    if require[line][0]["value"] != "":
		print Fore.GREEN+Style.BRIGHT+ "+ " +Style.RESET_ALL+line + ": " +value
	    else:
		print Fore.WHITE+Style.BRIGHT+ "* " +Style.RESET_ALL+line + "(" +Fore.GREEN+ "optional" +Style.RESET_ALL+ "):" +value
    #print Back.WHITE + Fore.WHITE + "End parameters" + Style.RESET_ALL


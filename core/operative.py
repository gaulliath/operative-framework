#!/usr/bin/env	python
# -*- coding: utf-8 -*-

import sys,os
import time
from core import menu
from core import mecanic
from colorama import Fore,Back,Style

def loading():
	action = 0
	base = "loading the fingerprinting framework"
	loading = base
	while action < 100:
		os.system('clear')
		if loading == base:
			loading = loading[0].upper() + loading[1:]
		print loading
		next = 0
		new_loading = ""
		for char in loading:
			if char.isupper():
				char = char.lower()
				next = 1
			elif next == 1 and char != " " and char != "*":
				char = char.upper()
				next = 0
			new_loading = new_loading + char
		if new_loading != "":
			loading = new_loading
		time.sleep(0.1)
		action += 1
	os.system('clear')
	return True


def user_put():
	version = "1.0a"
	print """                               __  _          
  ____  ____  ___  _________ _/ /_(_)   _____ 
 / __ \/ __ \/ _ \/ ___/ __ `/ __/ / | / / _ \\
/ /_/ / /_/ /  __/ /  / /_/ / /_/ /| |/ /  __/
\____/ .___/\___/_/   \__,_/\__/_/ |___/\___/ 
    /_/ """+Fore.RED+"Version: "+Style.RESET_ALL+version+" | "+Fore.RED+"Twitter: "+Style.RESET_ALL+"""@graniet75                               """
	action = 0
	print Fore.YELLOW + "        If you don't know how run it use :help" + Style.RESET_ALL
	print ""
	while action == 0:
		try:
			user_input = raw_input("$ operative > ")
		except:
			print "..."
			sys.exit()
		if ":" in user_input:
				user_input = user_input.replace(':','')
		if user_input in menu.menu_list:
			mecanic.load(menu.menu_list[user_input])
		if "use " in user_input:
			mecanic.load(user_input)
		

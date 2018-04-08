#!/usr/bin/env	python

from colorama import Fore,Back,Style
import os, sys
import requests
import re
from bs4 import BeautifulSoup

def show_options(require):
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


def export_data(export, export_file, export_status, title, argv):
	if len(export) > 0:
		if export_file == "":
			if argv == False:
				user_input = raw_input("operative (export file name ?) > ")
			else:
				user_input = argv
			if os.path.exists("export/"+user_input):
				export_file = "export/"+user_input
			elif os.path.exists(user_input):
				export_file = user_input
			else:
				print Fore.GREEN + "Writing " + user_input + " file" + Style.RESET_ALL
				export_file = "export/"+user_input
			export_data(export, export_file, export_status, title, argv)
		elif export_status == False:
			file_open = open(export_file, "a+")
			file_open.write(title)
			for line in export:
				try:
					file_open.write("- " + line +"\n")
				except:
					print Fore.RED + "ERROR: " + Style.RESET_ALL + "Can't write element."
			print Fore.GREEN + "File writed : " + export_file + Style.RESET_ALL
			file_open.close()
			export_status = True
		else:
			print Back.YELLOW + Fore.BLACK + "Module empty result" + Style.RESET_ALL

def export_data_search_db(export, export_file, export_status, title):
	if len(export) > 0:
		if export_file == "":
			user_input = raw_input("operative (export file name ?) > ")
		if os.path.exists("export/"+user_input):
			export_file = "export/"+user_input
		elif os.path.exists(user_input):
			export_file = user_input
		else:
			print Fore.GREEN + "Writing " + user_input + " file" + Style.RESET_ALL
		export_file = "export/"+user_input
		export_data(export, export_file, export_status, title)
	elif export_status == False:
		file_open = open(export_file,"a+")
		file_open.write(title)
		for line in export:
			file_open.write("- " + line +"\n")
		print Fore.GREEN + "File writed : " + export_file + Style.RESET_ALL
		file_open.close()
		export_status = True
	else:
		print Back.YELLOW + Fore.BLACK + "Module empty result" + Style.RESET_ALL

def set_options(require, name, value):
	if name in require:
		require[name][0]["value"] = value
	else:
		print Fore.RED + "Option not found" + Style.RESET_ALL

def check_require(require):
	for line in require:
		for option in require[line]:
			if option["required"] == "yes":
				if option["value"] == "":
					return False
	return True

def getDork(website, dorkName, browser):
	Valid = False
	if browser.strip().lower() == "google":
		server = "www.google.com"
	elif browser.strip().lower() == "bing":
		server = "www.bing.com"
	limit = 100
	url = "http://" + server + "/search?num=" + str(limit) + "&start=0&hl=en&meta=&q=" + str(dorkName)
	try:
		r = requests.get(url)
		valid = True
	except:
		valid = False
	if valid == True:
		print Fore.YELLOW + "Reading results of: " + str(dorkName) + Style.RESET_ALL
		result = r.content
		soup = BeautifulSoup(result, "html.parser")
		links = soup.findAll("a")
		if len(links) > 0:
			print Fore.YELLOW + "result: " + Style.RESET_ALL + str(len(links)) + " links found"
			correct_result = soup.find_all("a", href=re.compile("(?<=/url\?q=)(htt.*://.*)"))
			if len(correct_result) > 0:
				for link in soup.find_all("a", href=re.compile("(?<=/url\?q=)(htt.*://.*)")):
					print Fore.BLUE + " * " + Style.RESET_ALL +  re.split(":(?=http)", link["href"].replace("/url?q=", ""))[0]
				return correct_result
			else:
				print Fore.YELLOW + "result:" + Style.RESET_ALL + " No good links found"
				user_put = raw_input('$ operative (show '+str(len(links))+' other links? [N/y]) ')
				if user_put.lower() == "y":
					for link in links:
						print Fore.BLUE + " * " + Style.RESET_ALL + str(re.split(":(?=http)", link["href"].replace("/url?q=", ""))[0])
					return links
				else:
					return False
		else:
			print Fore.YELLOW + "result:" + Style.RESET_ALL + " No links found"
		return valid
	return valid
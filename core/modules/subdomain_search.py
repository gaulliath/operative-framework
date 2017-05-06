#!/usr/bin/env	python
#description:Search subdomain with google dork#

from colorama import Fore,Back,Style
from bs4 import BeautifulSoup

import os,sys
import requests
import urllib

class module_element(object):

	def __init__(self):
		self.title = "Subdomain : \n"
		self.require = {"website":[{"value":"","required":"yes"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False

	def set_agv(self, argv):
		self.argv = argv

	def show_options(self):
		#print Back.WHITE + Fore.WHITE + "Module parameters" + Style.RESET_ALL
		for line in self.require:
			if self.require[line][0]["value"] == "":
				value = "No value"
			else:
				value = self.require[line][0]["value"]
			if self.require[line][0]["required"] == "yes":
				if self.require[line][0]["value"] != "":
					print Fore.GREEN+Style.BRIGHT+ "+ " +Style.RESET_ALL+line+ ": " +value
				else:
					print Fore.RED+Style.BRIGHT+ "- " +Style.RESET_ALL+line+ "(" +Fore.RED+ "is_required" +Style.RESET_ALL+ "):" +value
			else:
				if self.require[line][0]["value"] != "":
					print Fore.GREEN+Style.BRIGHT+ "+ " +Style.RESET_ALL+line + ": " +value
				else:
					print Fore.WHITE+Style.BRIGHT+ "* " +Style.RESET_ALL+line + "(" +Fore.GREEN+ "optional" +Style.RESET_ALL+ "):" +value
		#print Back.WHITE + Fore.WHITE + "End parameters" + Style.RESET_ALL

	def export_data(self, argv=False):
		if len(self.export) > 0:
			if self.export_file == "":
				if argv == False:
					user_input = raw_input("operative (export file name ?) > ")
				else:
					user_input = argv
				if os.path.exists("export/"+user_input):
					self.export_file = "export/"+user_input
				elif os.path.exists(user_input):
					self.export_file = user_input
				else:
					print Fore.GREEN + "Writing " + user_input + " file" + Style.RESET_ALL
					self.export_file = "export/"+user_input
				self.export_data()
			elif self.export_status == False:
				file_open = open(self.export_file,"a+")
				file_open.write(self.title)
				for line in self.export:
					file_open.write("- " + line +"\n")
				print Fore.GREEN + "File writed : " + self.export_file + Style.RESET_ALL
				file_open.close()
				self.export_status = True
		else:
			print Back.YELLOW + Fore.BLACK + "Module empty result" + Style.RESET_ALL
	
	def set_options(self,name,value):
		if name in self.require:
			self.require[name][0]["value"] = value
		else:
			print Fore.RED + "Option not found" + Style.RESET_ALL
	
	def check_require(self):
		for line in self.require:
			for option in self.require[line]:
				if option["required"] == "yes":
					if option["value"] == "":
						return False
		return True

	def get_options(self,name):
		if name in self.require:
			return self.require[name][0]["value"]
		else:
			return False

	def run_module(self):
		ret = self.check_require()
		if ret == False:
			print Back.YELLOW + Fore.BLACK + "Please set the required parameters" + Style.RESET_ALL
		else:
			self.main()

	def main(self):
		server = "www.google.com"
		headers = {'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11',
       'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
       'Accept-Charset': 'ISO-8859-1,utf-8;q=0.7,*;q=0.3',
       'Accept-Encoding': 'none',
       'Accept-Language': 'en-US,en;q=0.8',
       'Connection': 'keep-alive'}
		website = self.get_options('website')
		if "://" in website:
			website = website.split('://')[1]
		dork = "site:"+website
		url="http://"+ server + "/search?num=100&start=0&hl=en&meta=&q=" + str(dork)
		print Fore.YELLOW +  Style.DIM + "[wait] try loading url..." + Style.RESET_ALL
		try:
			req = requests.get(url,headers=headers)
			html = req.content
		except:
			print Fore.RED + "[error] Can't load url. " + Style.RESET_ALL
		soup = BeautifulSoup(html,"html.parser",from_encoding="utf-8")
		parsing = soup.findAll('cite',attrs={'class':'_Rm'})
		for link in parsing:
			link = link.text.encode('utf-8')
			if "://" in link:
				link = link.split('://')[1]
			if "/" in link:
				link = link.split('/')[0]
			if " " in link:
				link = link.split(' ')[0]
			if website in link:
				if not link in self.export:
					self.export.append(link)
					print Fore.GREEN +"(FOUND) "+Style.RESET_ALL + str(link)






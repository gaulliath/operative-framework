#!/usr/bin/env	python
#Display header information#

from colorama import Fore,Back,Style
from core import load

import os,sys
import requests

class module_element(object):

	def __init__(self):
		self.title = "Header retrieval : \n"
		self.require = {"website":[{"value":"","required":"yes"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False

	def set_agv(self, argv):
		self.argv = argv

	def show_options(self):
                load.show_options(self.require)

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
		headers = {'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11',
		'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
		'Accept-Charset': 'ISO-8859-1,utf-8;q=0.7,*;q=0.3',
		'Accept-Encoding': 'none',
		'Accept-Language': 'en-US,en;q=0.8',
		'Connection': 'keep-alive'}
		website = self.get_options('website')
		loaded = 0
		if not "://" in website:
			website = "http://" + website
		print Fore.YELLOW + Style.DIM + "Try load '" + str(website) +"'"+ Style.RESET_ALL
		try:
			req = requests.get(website)
			loaded = 1
			print Fore.GREEN + Style.DIM + "[OK] successfully '" + str(website) +"'"+ Style.RESET_ALL
		except:
			print Fore.RED + "Can't load url '"+str(website)+"'" + Style.RESET_ALL
		if loaded == 1:
			print Fore.YELLOW + Style.DIM + "[wait] header recovery..." + Style.RESET_ALL
			headers = req.headers
			if len(headers) > 0:
				for line in headers:
					lines = line + " : " + headers[line]
					if not lines in self.export:
						print Fore.GREEN + "-" + Style.RESET_ALL + str(lines)
						self.export.append(lines)
			else:
				print Fore.RED + "Header is empty'"+str(website)+"'" + Style.RESET_ALL



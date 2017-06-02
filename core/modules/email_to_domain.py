#!/usr/bin/env	python
#description:Get domain with email#

from colorama import Fore,Back,Style
from core import load

import os,sys
import requests
import re
import string

class module_element(object):

	def __init__(self):
		self.title = "Email whois gathering : \n"
		self.require = {"email":[{"value":"","required":"yes"}]}
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
		headers = {
    		'User-Agent': 'Mozilla/5.0'
    	}
		email = self.get_options('email')
		url = "https://whoisology.com/search_ajax/search?action=email&value="+email+"&page=1&section=admin"
		output = ""
		try:
			output = requests.get(url,headers=headers)
			output = output.content
		except:
			print Fore.RED + "Can't open url" + Style.RESET_ALL
		if output != "":
			regex = re.compile('whoisology\.com\/(.*?)">')
			finded = regex.findall(output)
			if len(finded) > 0:
				for line in finded:
					if line.strip() != "":
						if line not in self.export and "." in line:
							self.export.append(line)
							print "- "+Fore.GREEN + line + Style.RESET_ALL
			else:
				print Fore.YELLOW + "Empty domain result for email : "+email +Style.RESET_ALL




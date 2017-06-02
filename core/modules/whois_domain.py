#!/usr/bin/env	python
#description:	Whois information for domain#

from colorama import Fore,Back,Style
from core import load

import os,sys
import pythonwhois
import urllib

class module_element(object):

	def __init__(self):
		self.title = "Whois domain gathering : \n"
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
		detail = None
		website = self.require["website"][0]["value"]
		if "://" in self.get_options('website'):
			website = self.get_options('website').split('://')[1]
		try:
			whois_information = pythonwhois.get_whois(website)
			detail = whois_information['contacts']['registrant']
		except:
			print Fore.RED + "Please use correct name without http(s)://" + Style.RESET_ALL
		export = []
		total = ""
		if detail != None:
			for element in detail:
				print Fore.BLUE + "* " + Style.RESET_ALL + element + " : " + str(detail[element])
				total = element + " : "+ str(detail[element])
				self.export.append(total)
		else:
			print Fore.RED + "Can't get whois information for : " +self.get_options('website') + Style.RESET_ALL

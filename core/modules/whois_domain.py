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
                load.export_data(self.export, self.export_file, self.export_status, self.title, argv)
	
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

#!/usr/bin/env	python
#description:Search default password from manufactor#

from colorama import Fore,Back,Style
from core import load

import os,sys
import requests
from bs4 import BeautifulSoup

class module_element(object):

	def __init__(self):
		self.title = "Default Password : \n"
		self.require = {"manufactor":[{"value":"","required":"yes"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False

	def set_agv(self, argv):
		self.argv = argv

	def show_options(self):
		return load.show_options(self.require)

	def export_data(self, argv=False):
		return load.export_data(self.export, self.export_file, self.export_status, self.title, argv)
	
	def set_options(self,name,value):
		return load.set_options(self.require, name, value)

	def check_require(self):
		return load.check_require(self.require)

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
		manufactor = str(self.get_options("manufactor"))
		search_link = "http://www.defaultpassword.com/?action=dpl&char="+manufactor
		print Fore.GREEN + "Searching for passwords" + Style.RESET_ALL
		req = requests.get(search_link)
		html = req.content
		result = html.split('<TR VALIGN="top">')
		try:
			result.pop(0)
			result.pop(0)
			print Back.BLACK + Fore.YELLOW + " ! Manufactor	|	Product		|	Protocol	|	User	|	Password" + Style.RESET_ALL
			for line in result:
				spliting = line.split('<TD NOWRAP>')
				string = ''
				for element in spliting:
					if "</TD>" in element and element.strip() != '':
						string = Fore.BLUE + " * " + Style.RESET_ALL + spliting[1].split('</TD>')[0] \
								 + Fore.YELLOW + "	|	" + Style.RESET_ALL + spliting[2].split('</TD>')[0] \
								 + Fore.YELLOW + "	|	" + Style.RESET_ALL + spliting[4].split('</TD>')[0] \
								 + Fore.YELLOW + "	|	" + Style.RESET_ALL + spliting[5].split('</TD>')[0] \
								 + Fore.YELLOW + "	|	" + Style.RESET_ALL + spliting[6].split('</TD>')[0]
						manufactor = spliting[1].split('</TD>')[0]
						product = spliting[2].split('</TD>')[0]
						protocol = spliting[4].split('</TD>')[0]
						user = spliting[5].split('</TD>')[0]
						password = spliting[6].split('</TD>')[0]
						string_export = {'manufactor':manufactor,'product':product,'protocol':protocol,'user':user,'password':password}
						self.export.append(string_export)
				if string != '':
					print string
		except IndexError:
			print Fore.RED + "No results were found for this manufacturer." + Style.RESET_ALL





#!/usr/bin/env	python
#description:Search possible SQL error index with google#

from colorama import Fore,Back,Style
from core import load

import os,sys
import urllib

class module_element(object):

	def __init__(self):
		self.title = "Check SQL error (Google Hacking) : \n"
		self.require = {"website": [{"value": "", "required": "yes"}]}
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
		websiteName = str(self.get_options('website').lower())
		if "://" in websiteName:
			websiteName = websiteName.split('://')[1]
		dorking = "site:" + str(websiteName) + ' intext:"sql syntax near" | intext:"syntax error has occurred" | intext:"incorrect syntax near" | intext:"unexpected end of SQL command" | intext:"Warning: mysql_connect()" | intext:"Warning: mysql_query()" | intext:"Warning: pg_connect()"'
		load.getDork(websiteName, dorking, browser="google")



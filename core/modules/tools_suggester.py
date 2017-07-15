#!/usr/bin/env	python
#description:Check website & show possible tools for CMS exploitation#

from colorama import Fore,Back,Style
from core import load

import os,sys
import requests

class module_element(object):

	def __init__(self):
		self.title = "CMS tools suggester : \n"
		self.require = {"website":[{"value":"","required":"yes"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False
		self.status_code = [200,403]
		self.tools = [
			{"tools":"wpscan","type":"wordpress","url":"https://github.com/wpscanteam/wpscan"},
			{"tools":"joomscan","type":"joomla","url":"https://sourceforge.net/projects/joomscan/"},
			{"tools":"drupscan","type":"drupal","url":"https://github.com/tibillys/drupscan"},
			{"tools":"SPIPScan","type":"spip","url":"https://github.com/PaulSec/SPIPScan"},
			{"tools":"Magescan","type":"magento","url":"https://github.com/steverobbins/magescan"}
		]
		self.directory = [
			{"file":"/wp-includes/","type":"wordpress","intext":""},
			{"file":"/wp-admin/","type":"wordpress","intext":""},
			{"file":"/readme.html","type":"wordpress","intext":"WordPress"},
			{"file":"/CHANGELOG.txt","type":"drupal","intext":"drupal"},
			{"file":"/administrator/","type":"joomla","intext":"Joomla"},
			{"file":"/robots.txt","type":"spip","intext":"SPIP"},
			{"file":"/frontend/default/","type":"magento","intext":""},
			{"file":"/static/frontend/","type":"magento","intext":""}
		]

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
		current = ""
		website = self.get_options('website')
		if website[-1] == "/":
			website = website[:-1]
		if "http" not in website or "https" not in website:
			website = "http://"+website
		for element in self.directory:
			complet = website + element['file']
			req = requests.get(complet)
			if req.status_code in self.status_code:
				if element["intext"] != "":
					if element["intext"].upper() in req.content.upper():
						current = element["type"]
					else:
						current = ""
				else:
					current = element["type"]
			if current != "":
				for tool in self.tools:
					if tool["type"].upper() == current.upper():
						if not tool["tools"] + " (" + tool["url"] + ")" in self.export:
							print Fore.GREEN + "Possible usage of " + str(tool["tools"].lower()) + " for " + str(current.lower()) + Style.RESET_ALL
							self.export.append(tool["tools"] + " (" + tool["url"] + ")")

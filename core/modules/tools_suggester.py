#!/usr/bin/env	python
#description:Check website & show possible tools for CMS exploitation#

from colorama import Fore,Back,Style

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
		#print Back.WHITE + Fore.WHITE + "Module parameters" + Style.RESET_ALL
		for line in self.require:
			if self.require[line][0]["value"] == "":
				value = "No value"
			else:
				value = self.require[line][0]["value"]
			if self.require[line][0]["required"] == "yes":
				print Fore.RED + Style.BRIGHT + "- "+Style.RESET_ALL + line + ":" + Fore.RED + "is_required" + Style.RESET_ALL + ":" + value
			else:
				print Fore.WHITE + Style.BRIGHT + "* "+Style.RESET_ALL + line + "(" + Fore.GREEN + "not_required" + Style.RESET_ALL + "):" + value
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
		current = ""
		website = self.get_options('website')
		if website[-1] == "/":
			website = website[:-1]
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

#!/usr/bin/env	python
#description:Reverse ip domain check (Yougetsignal)#

from colorama import Fore,Back,Style

import os,sys
import urllib
import requests
import json

class module_element(object):

	def __init__(self):
		self.title = "Reverse ip gathering : \n"
		self.require = {"domain":[{"value":"","required":"yes"}]}
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
	def is_array(self,var):
		return isinstance(var, (list, tuple))

	def main(self):
		content = ""
		url = "http://domains.yougetsignal.com/domains.php"
		try:
			r = requests.post(url, data = {'remoteAddress':self.get_options('domain')})
			content = r.json()
		except:
		 	print Fore.RED + "Can't send requests" + Style.RESET_ALL
		print Fore.GREEN + "Search information for : " + self.get_options('domain') + Style.RESET_ALL
		if content != "":
			for line in content:
				value = ""
				if self.is_array(content[line]):
					print "------------------------------"
					self.export.append("domain listing... : ")
					for content_array in content[line]:
						if self.is_array(content_array):
							value = "-----" + content_array[0]
							print Fore.BLUE + "* " + Style.RESET_ALL + content_array[0]
						else:
							value = "-----" + content_array
							print Fore.BLUE + "* " + Style.RESET_ALL + content_array
						self.export.append(value)
					print "------------------------------"
				else:
		 			print line + " : " + content[line]
		 			value = line + " : " + content[line]
		 		self.export.append(value)
		else:
			print Fore.YELLOW + "" + Style.RESET_ALL






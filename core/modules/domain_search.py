#!/usr/bin/env	python
#description:Search enterprise domain name#

from colorama import Fore,Back,Style

import os,sys
import urllib

class module_element(object):

	def __init__(self):
		self.title = "Domain gathering : \n"
		self.require = {"enterprise":[{"value":"","required":"yes"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False

	def set_agv(self, argv):
		self.argv = argv

	def show_options(self):
		#print Back.BLACK + Fore.WHITE + "==========" + Style.RESET_ALL
		for line in self.require:
			if self.require[line][0]["value"] == "":
				value = "No value"
			else:
				value = self.require[line][0]["value"]
			if self.require[line][0]["required"] == "yes":
				print Fore.RED + Style.BRIGHT + "- "+Style.RESET_ALL + line + ":" + Fore.RED + "is_required" + Style.RESET_ALL + ":" + value
			else:
				print Fore.WHITE + Style.BRIGHT + "* "+Style.RESET_ALL + line + "(" + Fore.GREEN + "not_required" + Style.RESET_ALL + "):" + value
		#print Back.WHITE + Fore.WHITE + "==========" + Style.RESET_ALL

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
        	domain_list = []
        	load_name = self.get_options("enterprise")
        	print Style.BRIGHT + Fore.BLUE + "Search domain name for "+load_name + Style.RESET_ALL
        	start_with = ["www.","http://","https://"]
        	end_with   = [".com",".fr",".org",".de",".eu"]
        	for line in start_with:
                	for end_line in end_with:
                        	domain = line + str(load_name) + end_line
                        	try:
                                	return_code = urllib.urlopen(domain).getcode()
                                	return_code = str(return_code)
                                	if return_code != "404":
                                        	domain_list.append(domain)
                                        	print Fore.GREEN + "- "+Style.RESET_ALL + domain
                        	except:
                                	Back.YELLOW + Fore.BLACK + "Can't get return code" + Style.RESET_ALL
		if len(domain_list) > 0:
			for domain in domain_list:
				self.export.append(domain)
		else:
			print Fore.RED + "No domain found" + Style.RESET_ALL



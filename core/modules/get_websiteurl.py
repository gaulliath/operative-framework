#!/usr/bin/env	python
#description:Extract url on website domain#

from colorama import Fore,Back,Style
from bs4 import BeautifulSoup

import os,sys
import time
import requests

class module_element(object):

	def __init__(self):
		self.title = "Url gathering : \n"
		self.require = {"website_url":[{"value":"","required":"yes"}],"page_limit":[{"value":"100","required":"no"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False
		self.already = []
		self.linked = []
		self.current_load = 1

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

	def parse_domain(self, url):
		if "://" in url:
			url = url.split('://')[1]
		if "." in url:
			url = url.split('.',1)[0]
		return url

	def extract_url(self, url):
		next_page = ""
		nexts = 0
		try:
			req = requests.get(url)
			nexts = 1
		except:
			print Fore.YELLOW + "* Can't open " + str(url) + Style.RESET_ALL
		if nexts == 1:
			if url[-1:] == "/":
				url = url[:-1]
			if url not in self.already and self.parse_domain(self.get_options('website_url')) in url:
				html = req.content
				soup = BeautifulSoup(html, "html.parser")
				link_count = len(soup.findAll('a'))
				print Fore.YELLOW + "* Load : " + str(url) + " with " + str(link_count) + " total link" + Style.RESET_ALL
				for a in soup.findAll('a'):
					try:
						if a['href'] != "":
							total_link = a['href']
							if total_link[:1] == "/":
								total_link = url + total_link
							elif total_link[:2] == "//":
								total_link = total_link.replace('//',url + "/")
							elif total_link[:1] == "#":
								total_link = url + "/" + total_link
							if "mailto:" not in total_link:
								if total_link not in self.export_file and a['href'] not in self.linked:
									if total_link != "":
										self.export.append(total_link)
								self.linked.append(a['href'])
					except:
						print Fore.RED + "Can't read link" + Style.RESET_ALL
				if self.current_load <= int(self.get_options('page_limit')):
					if len(self.export) > 0:
						next_page = self.export[self.current_load]
						if next_page != "":
							self.current_load += 1
							self.extract_url(next_page)
				else:
					return True
	def main(self):
		nexts = 0
		website = self.get_options('website_url')
		if "http//" in website:
			website = website.replace('http//','http://')
		print Fore.GREEN + "* Check if " + str(website) + " is stable" + Style.RESET_ALL
		try:
			req = requests.get(website)
			nexts = 1
		except:
			print Fore.RED + "* Website url is not stable" + Style.RESET_ALL
		if nexts == 1:
			self.extract_url(website)


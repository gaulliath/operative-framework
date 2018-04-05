#!/usr/bin/env	python
#description:Search archive of website domain (archive.org)#

from colorama import Fore,Back,Style
from core import load

import os,sys
import json
import datetime
import requests

class module_element(object):

	def __init__(self):
		self.title = "Archive.org Gathering : \n"
		self.require = {"domain":[{"value":"","required":"yes"}],"from_date":[{"value":"2010","required":"no"}],"to_date":[{"value":"2017","required":"no"}],"limit":[{"value":"100","required":"no"}]}
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
		error = 0
		domain_name = str(self.get_options('domain'))
		from_date = str(self.get_options('from_date'))
		to_date = str(self.get_options('to_date'))
		limit_result = str(self.get_options('limit'))
		if domain_name[-1] == '/':
			domain_name = domain_name[:-1]
		if "://" in domain_name:
			domain_name = domain_name.split('://')[1]
		url = "http://web.archive.org/cdx/search/cdx?url="+domain_name+"&matchType=domain&limit="+limit_result+"&output=json&from="+from_date+"&to="+to_date
		try:
			req = requests.get(url)
			json_data = json.loads(req.text)
			if len(json_data) == 0:
				print Fore.YELLOW + "output: " + Style.RESET_ALL + "No result found"
				self.export.append("no result in archive")
				error = 1
		except:
			print Fore.RED + "error: " + Style.RESET_ALL + " Can't open url"
			error = 1
		if error == 0:
			try:
				result = [ x for x in json_data if x[2] != 'original']
				result.sort(key=lambda x: x[1])
				for line in result:
					timestamp = line[1]
					website   = line[2]
					total_link  = "https://web.archive.org/web/" + str(timestamp) + "/" + str(website)
					string_date = str(timestamp[:4]) + "/" + str(timestamp[4:6]) + "/" + str(timestamp[6:8])
					self.export.append(total_link)
					print " {}   {}  {}({}{}{}){}".format(string_date, website, Fore.YELLOW, \
						Style.RESET_ALL, total_link, Fore.YELLOW, Style.RESET_ALL)
			except:
				print Fore.RED + "error: " + Style.RESET_ALL + "Error found please retry"
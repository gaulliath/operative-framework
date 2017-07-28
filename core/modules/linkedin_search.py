#!/usr/bin/env	python
#description:Linkedin employee search module#

from colorama import Fore,Back,Style
from bs4 import BeautifulSoup
from core import load

import os,sys
import requests

class module_element(object):

	def __init__(self):
		self.title = "Linkedin gathering : \n"
		self.require = {"enterprise":[{"value":"","required":"yes"}],"limit_search":[{"value":"","required":"yes"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False

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

	def set_agv(self, argv):
		self.argv = argv

	def run_module(self):
		ret = self.check_require()
		if ret == False:
			print Back.YELLOW + Fore.BLACK + "Please set the required parameters" + Style.RESET_ALL
		else:
			self.main()

	def main(self):
            server = "encrypted.google.com"
            limit = self.get_options("limit_search")
            enterprise = self.get_options("enterprise")
            counter = ""
            url = "https://" + server + "/search?num=" + limit + "&start=0&hl=en&q=site:linkedin.com/in+" + enterprise
            print Fore.GREEN + "Search Linkedin research" + Style.RESET_ALL
            request = requests.get(url)
	    status_code = request.status_code

	    if status_code == 200:
                html = BeautifulSoup(request.text, "html.parser")
        	results = html.find_all('div', { 'class' : 'g' })
		
		for i, result in enumerate(results):
                    employee = result.find('h3', { 'class' : 'r' }).getText()
		    if "| LinkedIn" or "on LinkedIn" or "LinkedIn" in employee:
		        employee = employee.replace('| LinkedIn', '')
		        employee = employee.replace('LinkedIn', '')
		        employee = employee.replace('on LinkedIn', '')
		
		    profile = result.find('cite').getText()
                    print Fore.BLUE + "* "+ Style.RESET_ALL + employee + " < " + Fore.GREEN + profile + Style.RESET_ALL
                    counter = i

                if counter == "":
                    print Fore.RED + "Nothing on linkedin." + Style.RESET_ALL
                else:
                    print "\n Total results:", counter+1
	    else:
                print Fore.RED + "Can't get response" + Style.RESET_ALL

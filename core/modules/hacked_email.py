#!/usr/bin/env	python
# description:Module sample#

from colorama import Fore, Back, Style
from core import load

import os, sys
import urllib
import requests
import json


class module_element(object):
    def __init__(self):
        self.title = "Hacked email : \n"
        self.require = {"email": [{"value": "", "required": "yes"}]}
        self.export = []
        self.export_file = ""
        self.export_status = False

    def set_agv(self, argv):
        self.argv = argv

    def show_options(self):
        return load.show_options(self.require)

    def export_data(self, argv=False):
        return load.export_data(self.export, self.export_file, self.export_status, self.title, argv)

    def set_options(self, name, value):
        return load.set_options(self.require, name, value)

    def check_require(self):
        return load.check_require(self.require)

    def get_options(self, name):
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
        api_url = 'https://hacked-emails.com/api?q='
        email = self.get_options('email')
        if email != "":
            if "@" in email and "." in email:
                complet_url = api_url + str(email)
                req = requests.get(complet_url)
                content = req.text
                if content != "":
                    content = json.loads(content)
                    if content['status'] and content['status'] == "found":
                        print "Result found (" + Fore.GREEN + str(content['results']) + " results" + Style.RESET_ALL + ")"
                        for line in content['data']:
                            try:
                                print Fore.BLUE + " * " + Style.RESET_ALL + " found in : " + Fore.GREEN + str(line['title']) + Style.RESET_ALL + \
                                      " (" + Fore.YELLOW + str(line['date_leaked']) + Style.RESET_ALL + ")"
                            except:
                                print Fore.BLUE + " * " + Style.RESET_ALL + " found in : ( can't parse leaks title)"
                    else:
                        print "Status (" + Fore.RED + "Not found" + Style.RESET_ALL + ")"
                else:
                    print Fore.RED + "Error found in json" + Style.RESET_ALL + ")"
            else:
                print Fore.YELLOW + "Invalid email please retry with correct email address" + Style.RESET_ALL



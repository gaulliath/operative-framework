#!/usr/bin/env  python
#description:Check syntax for operative framework modules#
# -*- coding: utf-8 -*-

import os,sys
from colorama import Fore, Style, Back

class helper_class(object):

    def __init__(self):
        self.resume  = {
            'Name': 'Module Class Syntax checker',
            'Desc': 'Check syntax for operative framework modules',
            'Author': 'Tristan Granier'
        }
        self.require = {"module_path": [{"value": "", "required": "yes"}]}

    def set_agv(self, argv):
        self.argv = argv

    def get_resume(self):
        print Fore.YELLOW + "Name          " + Style.RESET_ALL + ": " + self.resume['Name']
        print Fore.YELLOW + "Description   " + Style.RESET_ALL + ": " + self.resume['Desc']
        print Fore.YELLOW + "Author        " + Style.RESET_ALL + ": (c) " + self.resume['Author']

    def show_options(self):
        # print Back.WHITE + Fore.WHITE + "Module parameters" + Style.RESET_ALL
        for line in self.require:
            if self.require[line][0]["value"] == "":
                value = "No value"
            else:
                value = self.require[line][0]["value"]
            if self.require[line][0]["required"] == "yes":
                if self.require[line][0]["value"] != "":
                    print Fore.GREEN + Style.BRIGHT + "+ " + Style.RESET_ALL + line + ": " + value
                else:
                    print Fore.RED + Style.BRIGHT + "- " + Style.RESET_ALL + line + "(" + Fore.RED + "is_required" + Style.RESET_ALL + "):" + value
            else:
                if self.require[line][0]["value"] != "":
                    print Fore.GREEN + Style.BRIGHT + "+ " + Style.RESET_ALL + line + ": " + value
                else:
                    print Fore.WHITE + Style.BRIGHT + "* " + Style.RESET_ALL + line + "(" + Fore.GREEN + "optional" + Style.RESET_ALL + "):" + value
                    # print Back.WHITE + Fore.WHITE + "End parameters" + Style.RESET_ALL


    def set_options(self, name, value):
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

    def reload_main(self):
        user_put = raw_input('(operative) Reload module ? [Y/n]')
        if user_put == "" or user_put == "Y" or user_put == "y":
            self.run_module()

    def main(self):
        module_name = self.get_options('module_path')
        success = 0
        if os.path.isfile(module_name):
            if "/" in module_name:
                module_path = module_name.replace("/", ".")
                module_path = module_path.replace('.py','')
                try:
                    mod = __import__(module_path, fromlist=['module_element'])
                    module_class = mod.module_element()
                    success = 1
                    if success == 1:
                        print "------------------------"
                        try:
                            success = 0
                            print Fore.YELLOW + "Information: " + Style.RESET_ALL + " :show options"
                            module_class.show_options()
                            success = 1
                            print "+ status: " + Fore.GREEN + "OK" + Style.RESET_ALL
                            if success == 1:
                                try:
                                    success = 0
                                    print Fore.YELLOW + "Information: " + Style.RESET_ALL + " :set options=value"
                                    for line in module_class.require:
                                        print "- set " + str(line)
                                        module_class.set_options(line,'val_debug')
                                    print "+ status: " + Fore.GREEN + "OK" + Style.RESET_ALL
                                    success = 1
                                except Exception, e:
                                    print "- status: " + Fore.RED + "ERROR" + Style.RESET_ALL
                                    print str(e)
                                    self.reload_main()
                            if success == 1:
                                try:
                                    success = 0
                                    print Fore.YELLOW + "Information: " + Style.RESET_ALL + " :run"
                                    module_class.run_module()
                                    print "+ status: " + Fore.GREEN + "OK" + Style.RESET_ALL
                                    success = 1
                                except Exception,e:
                                    print "- status: " + Fore.RED + "ERROR" + Style.RESET_ALL
                                    print str(e)
                                    self.reload_main()
                                if success == 1:
                                    print Fore.GREEN + "Your module seem be OK" + Style.RESET_ALL
                                    print "------------------------"


                        except Exception,e:
                            print "- status: " + Fore.RED + "ERROR" + Style.RESET_ALL
                            print str(e)
                            self.reload_main()
                except Exception,e:
                    print "- status: " + Fore.RED + "ERROR" + Style.RESET_ALL
                    print str(e)
                    self.reload_main()
            else:
                Fore.RED + "- " + Style.RESET_ALL + "Please use full relative path of module"
        else:
            print Fore.RED + "- " + Style.RESET_ALL + "Please enter correct module path"

#!/usr/bin/env	python
# -*- coding: utf-8 -*-

import sys,os
import time
import glob
import urllib
import string
import random
import json
from colorama import Fore,Back,Style
from core import export

def SetError(value):
    return {'error':value}

def API_listmodule(argument):
    total_module = {'module_list':[]}
    if os.path.exists("core/modules/"):
        list_module = glob.glob("core/modules/*.py")
        for module in list_module:
            if ".py" in module:
                module_name = module.split(".py")[0]
                module_name = module_name.replace('core/modules/', '')
            if "__init__" not in module:
                description = "No module description found"
                if "#description:" in open(module).read():
                    description = open(module).read().split("#description:")[1]
                    description = description.split("#")[0]
                    total_module['module_list'].append({'name':module_name,'description':description})
    if len(total_module['module_list']) > 0:
        return total_module
    else:
        return False

def API_startmodule(arguments):
    module_path = ''
    if len(arguments) > 0:
        if 'module_name' in arguments:
            if 'core/' in arguments['module_name']:
                if os.path.isfile(arguments['module_name']):
                    module_path = arguments['module_name'] + '.py'
                else:
                    return SetError("can't find module path")
            else:
                full_path = 'core/modules/'+arguments['module_name']+'.py'
                if os.path.isfile(full_path):
                    module_path = 'core.modules.'+arguments['module_name'] + '.py'
                else:
                    return SetError("can't find module path")
            if module_path != '':
                mod = __import__(module_path.replace("/", ".").split('.py')[0], fromlist=['module_element'])
                module_class = mod.module_element()
                for line in module_class.require:
                    if str(module_class.require[line][0]['required']) == 'yes':
                        if line in arguments:
                            module_class.set_options(line, arguments[line])
                        else:
                            print "can't set arguments : " + str(line)
                            return SetError({'message': "can't set arguments : " + str(line), 'fields': module_class.require})
                module_class.run_module()
                print module_class.export
                if len(module_class.export) > 0:
                    return module_class.export
                else:
                    return 'Module executed but nothing on export'
            else:
                print "can't find module path"
                return SetError("can't find module path")
        else:
            print "nop"
    else:
        return False

def API_requirementModule(arguments):
    module_path = ''
    if len(arguments) > 0:
        if 'module_name' in arguments:
            if 'core/' in arguments['module_name']:
                if os.path.isfile(arguments['module_name']):
                    module_path = arguments['module_name'] + '.py'
                else:
                    return SetError("can't find module path")
            else:
                full_path = 'core/modules/' + arguments['module_name'] + '.py'
                if os.path.isfile(full_path):
                    module_path = 'core.modules.' + arguments['module_name'] + '.py'
                else:
                    return SetError("can't find module path")
            if module_path != '':
                mod = __import__(module_path.replace("/", ".").split('.py')[0], fromlist=['module_element'])
                module_class = mod.module_element()
                return module_class.require
            else:
                print "can't find module path"
                return SetError("can't find module path")
        else:
            print "nop"
    else:
        return False

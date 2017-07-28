#!/usr/bin/env	python
# -*- coding: utf-8 -*-

import sys,os
import subprocess
import json
from threading import Thread
from colorama import Fore,Back,Style
from core import api

def loadwebserver():
    print Fore.YELLOW + '* For commercial/enterprise use please send email to graniet75@gmail.com' + Style.RESET_ALL
    cmd = "php"
    if not os.path.isfile('config.json'):
        print Fore.RED + "Config file not found" + Style.RESET_ALL
        sys.exit()
    config_file = open('config.json').read()
    port = int(json.loads(config_file)['webserver']['interface'])
    port_api = int(json.loads(config_file)['webserver']['api_port'])
    open_tpl_javascript = open('core/interface/operative_tpl.js').read()
    open_tpl_javascript = open_tpl_javascript.replace('{{PORT_API}}', str(port_api))
    open_file = open('core/interface/operative.js','w+')
    open_file.write(open_tpl_javascript)
    open_file.close()
    background_thread = ''
    try:
        FNULL = open(os.devnull, 'w')
        try:
            background_thread = Thread(target=api.loadwebserverapi, args=['NO_DEBUG'])
            background_thread.daemon = True
            background_thread.start()
            print Fore.BLUE + "*" + Style.RESET_ALL + " API service running to : http://127.0.0.1:" + str(port_api)
        except:
            print Fore.RED + "Can't run webserver api" + Style.RESET_ALL
            sys.exit()
        print Fore.BLUE + "*" + Style.RESET_ALL + " Web interface running to : http://127.0.0.1:"+str(port)
        subprocess.call([cmd, '-S', '127.0.0.1:'+str(port), '-t', 'core/interface/'], stdout=FNULL, stderr=subprocess.STDOUT)
    except:
        print "please install php cli"



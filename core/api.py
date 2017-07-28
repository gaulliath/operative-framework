#!/usr/bin/env	python
# -*- coding: utf-8 -*-
import sys,os
import cgi
import urlparse
import json
from BaseHTTPServer import BaseHTTPRequestHandler
from BaseHTTPServer import HTTPServer
from colorama import Fore,Back,Style
from core import api_shortcode
from core import mecanic_api

ERROR_FIELD = ''

def run_api(command, arguments):
    if command in api_shortcode.menu_list_api:
        result = getattr(mecanic_api, api_shortcode.menu_list_api[command])(arguments)
        if result != False:
            return result
        else:
            return False
    else:
        print False

def loadwebserverapi(debug='DEBUG'):
    config_file = open('config.json').read()
    port = int(json.loads(config_file)['webserver']['api_port'])
    server = HTTPServer(('localhost', port), GetHandler)
    if debug == 'DEBUG':
        print Fore.YELLOW + '* For commercial/enterprise use please send email to graniet75@gmail.com' + Style.RESET_ALL
        print Fore.BLUE + '* ' + Style.RESET_ALL + 'Starting server to http://127.0.0.1:'+str(port)
        print Fore.BLUE + '* ' + Style.RESET_ALL + 'use <Ctrl-C> to stop'
    server.serve_forever()

class GetHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_path = urlparse.urlparse(self.path)
        message_parts = [
            'CLIENT VALUES:',
            'client_address=%s (%s)' % (self.client_address,
                                        self.address_string()),
            'command=%s' % self.command,
            'path=%s' % self.path,
            'real path=%s' % parsed_path.path,
            'query=%s' % parsed_path.query,
            'request_version=%s' % self.request_version,
            '',
            'SERVER VALUES:',
            'server_version=%s' % self.server_version,
            'sys_version=%s' % self.sys_version,
            'protocol_version=%s' % self.protocol_version,
            '',
            'HEADERS RECEIVED:',
        ]
        for name, value in sorted(self.headers.items()):
            message_parts.append('%s=%s' % (name, value.rstrip()))
        message_parts.append('')
        return_json = {
            'status':'INFORMATION',
            'message':'Welcome to operative framework API',
            'data':'Please don\'t use (GET) but (POST)',
            'error':''
        }
        message = '\r\n'.join(message_parts)
        self.send_response(200)
        self.end_headers()
        self.wfile.write(json.dumps(return_json))
        return True

    def do_POST(self):
        form = cgi.FieldStorage(
            fp=self.rfile,
            headers=self.headers,
            environ={'REQUEST_METHOD': 'POST',
                     'CONTENT_TYPE': self.headers['Content-Type'],
                     })
        self.send_response(200)
        self.send_header('Content-Type','json')
        self.send_header("Access-Control-Allow-Origin", "*")
        self.end_headers()
        current_command = ''
        current_value = {}
        message = 'command correctly executed'
        for field in form.keys():
            field_item = form[field]
            if field == 'exec':
                status = 'OK'
                error = ''
                current_command = form[field].value
            else:
                current_value[field] = form[field].value
            current_command = current_command.lower()
        result = run_api(current_command, current_value)
        if not result:
            result = ""
            message = 'Error has occured'
            status = 'ERROR'
            error = 'Unknow error'
        elif 'error' in result:
            message = 'Error has occured'
            status = 'ERROR'
            error = result['error']
            result = ""
        json_format = {
            'status':status,
            'message':message,
            'data':result,
            'error':error
        }
        self.wfile.write(json.dumps(json_format, ensure_ascii=False))
        return True

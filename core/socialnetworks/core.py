#!/usr/bin/env	python
# -*- coding: utf-8 -*-

import os
import sys
import json
import requests
from colorama import Fore,Back,Style


class SocialNetwork(object):

    def __init__(self):
        self.type_of_search = [
            {'name': 'twitter', 'resume': 'Search '},
        ]

    def run(self):
        print "This module use " + Fore.YELLOW + "fullcontact.com" + Style.RESET_ALL
        print Fore.BLUE + "* " + Style.RESET_ALL + "checking configuration..."
        if os.path.isfile('config.json'):
            fullcontact_api = False
            with open('config.json') as json_file:
                data_json = json.load(json_file)
            try:
                fullcontact_api = data_json['external_api']['fullcontact']
            except:
                print Fore.RED + "Please configure fullcontact API is json" + Style.RESET_ALL
            if fullcontact_api != False:
                action = 1
                email_input = False
                while action == 1:
                    email_input = raw_input('$ operative ('+ Fore.YELLOW +'Social Network'+ Style.RESET_ALL +') : please enter email like jhon@doe.com: ')
                    if email_input != '':
                        if '@' in email_input and '.' in email_input:
                            action = 0
                        else:
                            print Fore.YELLOW + "error:" + Style.RESET_ALL + " Please enter correct email address."
                headers = {
                    'X-FullContact-APIKey': fullcontact_api,
                }

                params = (
                    ('email', email_input),
                )

                content = requests.get('https://api.fullcontact.com/v2/person.json', headers=headers, params=params)
                if content:
                    json_response = json.loads(content.text)
                    if json_response['status'] == 200:
                        social_network = []
                        contact_info = False
                        localization = False
                        photos = []
                        result_count = 0
                        try:
                            social_network = json_response['socialProfiles']
                            result_count += 1
                        except:
                            print Fore.RED + "* " + Style.RESET_ALL + " Unknown contact social network"

                        try:
                            contact_info = json_response['contactInfo']
                            result_count += 1
                        except:
                            print Fore.RED + "* " + Style.RESET_ALL + " Unknown contact information"

                        try:
                            localization  = json_response['demographics']
                            result_count += 1
                        except:
                            print Fore.RED + "* " + Style.RESET_ALL + " Unknown contact localization"

                        try:
                            photos = json_response['photos']
                            result_count += 1
                        except:
                            print Fore.RED + "* " + Style.RESET_ALL + " Unknown contact photos"

                        print "Result found ("+str(result_count)+"):"

                        if len(social_network) > 0:
                            for line in social_network:
                                print Fore.YELLOW + line['type'] + Style.RESET_ALL + ' : ' + line['url']
                            print "-----"
                        if contact_info != False:
                            if 'familyName' in contact_info:
                                print Fore.YELLOW + "Lastname : " + Style.RESET_ALL + contact_info['familyName']
                            if 'givenName' in contact_info:
                                print Fore.YELLOW + "FirstName: " + Style.RESET_ALL + contact_info['givenName']
                            if 'fullName' in contact_info:
                                print Fore.YELLOW + "FullName : " + Style.RESET_ALL + contact_info['fullName']
                            print "-----"
                        if localization != False:
                            if localization['locationGeneral']:
                                print Fore.YELLOW + "Localization: " + Style.RESET_ALL + localization['locationGeneral']
                                print "-----"
                        if len(photos) > 0:
                            for line in photos:
                                print Fore.YELLOW + "picture:" + Style.RESET_ALL + line['url'] + Fore.BLUE + ' <' + line['typeName'] + '>' + Style.RESET_ALL
                    if json_response['status'] == 403:
                        print Fore.RED + "error: " + Style.RESET_ALL + "Bad API key please get it https://fullcontact.com."
                    else:
                        print Fore.YELLOW + "Nothing returned..." + Style.RESET_ALL
                else:
                    print Fore.RED + "error: " + Style.RESET_ALL + "Can't return results..."
                    print Fore.BLUE+ "information: " + Style.RESET_ALL + "Are you sure your API key is right (fullcontact.com)?"

        else:
            print Fore.RED + "Can't locate a config.json" + Style.RESET_ALL



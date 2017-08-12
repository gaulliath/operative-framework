#!/usr/bin/env	python
# -*- coding: utf-8 -*-

import os,sys
import requests


class Filters(object):

    def __init__(self):
        # work_with variable, is used for check if current module
        # can be used with this filters
        # e.g: self.work_with=['linkedin_search', '...']
        # e.g: self.work_with = ['*'] for all module
        self.work_with = [
            'linkedin_search'
        ]

    # This function as started with result exported in current module
    def run(self, data):
        # Filter run after all module result
        # This filter show only employee name without link for the module linkedin_search
        if(len(data) > 0):
            for line in data:
                print line['employee']
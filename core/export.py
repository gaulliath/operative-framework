#!/usr/bin/env	python

import os,sys
from core import menu
from colorama import Fore,Back,Style

# global function
def begin(export_type,export_name):
	if export_type.upper() in menu.menu_export:
		try:
			globals()[menu.menu_export[export_type.upper()]["BEGIN"]](export_name)
			return True
		except:
			return False
	return False

def to(export_type,module_name,report_array,output_name,total_report=False):
	if export_type.upper() in menu.menu_export:
		globals()[menu.menu_export[export_type.upper()]["EXPORT"]](module_name,report_array,output_name,total_report)
		return True
	else:
		return False

def end(export_type,export_name):
	if export_type.upper() in menu.menu_export:
		try:
			globals()[menu.menu_export[export_type.upper()]["END"]](export_name)
			return True
		except:
			return False
	return False
# global function

# XML module function
def begin_module_XML(export_name):
	export_name = export_name.replace('.txt','.xml')
	file_open = open("export/" + export_name,'a+')
	file_open.write('<?xml version="1.0" encoding="UTF-8"?>\n')
	file_open.write('<operative-framework-report>\n')
	file_open.close()
	return True

def end_module_XML(export_name):
	export_name = export_name.replace('.txt','.xml')
	file_open = open("export/" + export_name,'a+')
	file_open.write('</operative-framework-report>')
	file_open.close()

def export_module_XML(export_name,export_array,output_name,total_report=False):
	output_name = output_name.replace('.txt','.xml')
	first_open = 0
	if len(export_array) > 0:
		nb = 1
		if ":" in export_name:
			export_name= export_name.replace(':', '')
		if '(' in export_name:
			export_name = export_name.replace('(','')
			export_name = export_name.replace(')','')
		export_name = export_name.strip()
		if " " in export_name:
			export_name = export_name.replace(' ', '-')
		# export_name_first = "<" + export_name + ">"
		# export_name_end = "</" + export_name + ">"
		export_name_first = "<report"+str(total_report)+">"
		export_name_end = "</report"+str(total_report)+">"
		file_open = open("export/"+output_name,'a+')
		file_open.write(export_name_first+"\n")
		file_open.write("	<name>"+export_name+"</name>\n")
		file_open.write("	<count>"+str(len(export_array))+"</count>\n")
		for line in export_array:
			if "-" in line[0]:
				line = line[0].replace('-','')
			line = "<value"+str(nb)+">"+line.strip()+"</value"+str(nb)+">"
			file_open.write("	"+line+"\n")
			nb+=1
		file_open.write(export_name_end +"\n")
# XML module function
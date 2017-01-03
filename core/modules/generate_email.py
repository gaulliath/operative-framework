#!/usr/bin/env	python
#description:Generate email with employee list#

from colorama import Fore,Back,Style

import os,sys
import urllib

class module_element(object):

	def __init__(self):
		self.title = "Email generator : \n"
		self.require = {"filename":[{"value":"","required":"yes"}]}
		self.export = []
		self.export_file = ""
		self.export_status = False
		self.domain = [
						'@gmail.com','@hotmail.com','@yahoo.com','@hotmail.fr',
						'@yahoo.fr','yandex.ru']

	def set_agv(self, argv):
		self.argv = argv

	def show_options(self):
		#print Back.WHITE + Fore.WHITE + "Module parameters" + Style.RESET_ALL
		for line in self.require:
			if self.require[line][0]["value"] == "":
				value = "No value"
			else:
				value = self.require[line][0]["value"]
			if self.require[line][0]["required"] == "yes":
				print Fore.RED + Style.BRIGHT + "- "+Style.RESET_ALL + line + ":" + Fore.RED + "is_required" + Style.RESET_ALL + ":" + value
			else:
				print Fore.WHITE + Style.BRIGHT + "* "+Style.RESET_ALL + line + "(" + Fore.GREEN + "not_required" + Style.RESET_ALL + "):" + value
		#print Back.WHITE + Fore.WHITE + "End parameters" + Style.RESET_ALL

	def export_data(self, argv=False):
		if len(self.export) > 0:
			if self.export_file == "":
				if argv == False:
					user_input = raw_input("operative (export file name ?) > ")
				else:
					user_input = argv
				if os.path.exists("export/"+user_input):
					self.export_file = "export/"+user_input
				elif os.path.exists(user_input):
					self.export_file = user_input
				else:
					print Fore.GREEN + "Writing " + user_input + " file" + Style.RESET_ALL
					self.export_file = "export/"+user_input
				self.export_data()
			elif self.export_status == False:
				file_open = open(self.export_file,"a+")
				file_open.write(self.title)
				for line in self.export:
					file_open.write("- " + line +"\n")
				print Fore.GREEN + "File writed : " + self.export_file + Style.RESET_ALL
				file_open.close()
				self.export_status = True
		else:
			print Back.YELLOW + Fore.BLACK + "Module empty result" + Style.RESET_ALL
	
	def set_options(self,name,value):
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

	def generate_email(self, name_last):
		email_list = []
		for domain in self.domain:
			email_first = name_last[1]+"."+name_last[0]+domain
			email_list.append(email_first)
			email_second= name_last[0]+"."+name_last[1]+domain
			email_list.append(email_second)
			email_third = name_last[1]+"-"+name_last[0]+domain
			email_list.append(email_third)
			email_four  = name_last[0]+"-"+name_last[1]+domain
			email_list.append(email_four)
			email_five  = name_last[0]+name_last[1]+domain
			email_list.append(email_five)
			email_six   = name_last[1]+name_last[0]+domain
			email_list.append(email_six)
		for email in email_list:
			self.export.append(email)


	def main(self):
		view = 0
		if os.path.exists(self.get_options('filename')):
			file_open = open(self.get_options('filename')).read()
			if "Viadeo gathering :" in file_open:
				print Fore.GREEN + "* "+Style.RESET_ALL + "Viadeo find..."
				view = 1
				explode_viadeo = file_open.split('Viadeo gathering :')[1]
				explode_viadeo = explode_viadeo.split('\n')
				for employee in explode_viadeo:
					if "-" in employee:
						employee = employee.split('-')[1].strip()
						if "." in employee:
							name_last = employee.split('.')
							self.generate_email(name_last)
			if "Linkedin gathering :" in file_open:
				print Fore.GREEN + "* "+Style.RESET_ALL + "Linkedin find..."
				view = 1
				explode_linkedin = file_open.split('Linkedin gathering :')[1]
				explode_linkedin = explode_linkedin.split('-')
				for employee in explode_linkedin:
					if "-" in employee:
						employee = employee.split('-')
						for line in employee:
							if " " in employee:
								name_last = employee.split(' ')
								self.generate_email(name_last)
				if view == 0:
					print Fore.YELLOW + "Please run Linkedin,Viadeo search module and export"+Style.RESET_ALL
				elif view==1:
					print Fore.GREEN + "* "+Style.RESET_ALL + "All email generated please use :export"



		else:
			print Fore.RED + self.get_options('filename')+Style.RESET_ALL+" Not valid file"



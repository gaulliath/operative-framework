#### Api usage 
#####  version.BETA

this API is current BETA version
Operative framework API accept POST method and return JSON data of module

available command:

    + list_module

This return list of all modules in json data e.g:

command : CURL -X POST 'http://127.0.0.1:9090/' -d '&exec=list_module'

```json
{"status": "OK", "message": "command correctly executed", "data": {"module_list": [{"name": "cms_gathering", "description": "Check if CMS is used (wordpress,joomla,magento)"}, {"name": "domain_search", "description": "Search enterprise domain name"}, {"name": "email_to_domain", "description": "Get domain with email"}, {"name": "file_common", "description": "\tRead/Search common file"}, {"name": "generate_email", "description": "Generate email with employee list"}, {"name": "get_websiteurl", "description": "Extract url on website domain"}, {"name": "getform_data", "description": "\tGet all form parameters (BETA)"}, {"name": "https_gathering", "description": "SSL/TLS information gathering (sslyze)"}, {"name": "linkedin_search", "description": "Linkedin employee search module"}, {"name": "metatag_look", "description": "\tget meta name,content"}, {"name": "reverse_ipdomain", "description": "Reverse ip domain check (Yougetsignal)"}, {"name": "sample_module", "description": "Module sample"}, {"name": "search_db", "description": "\tForensics module for SQL database"}, {"name": "subdomain_search", "description": "Search subdomain with google dork"}, {"name": "tools_suggester", "description": "Check website & show possible tools for CMS exploitation"}, {"name": "vhost_IPchecker", "description": "Reverse IP domain check (BING)"}, {"name": "viadeo_search", "description": "Viadeo employee search module"}, {"name": "waf_gathering", "description": "WAF information gathering : need wafw00f"}, {"name": "whois_domain", "description": "\tWhois information for domain"}]}, "error": ""}
```

    + requirement_module
This return all options needed from module name e.g:

command: CURL -X POST 'http://127.0.0.1:9090/' -d '&exec=requirement_module&module_name=linkedin_search'

```json
{"status": "OK", "message": "command correctly executed", "data": {"limit_search": [{"required": "yes", "value": ""}], "enterprise": [{"required": "yes", "value": ""}]}, "error": ""}
```

    + use_module

This execute current module with argument in POST element e.g:

command: CURL -X POST 'http://127.0.0.1:9090/' -d '&exec=use_module&module_name=linkedin_search&limit_search=10&enterprise=lynxframework'

```json
{"status": "OK", "message": "command correctly executed", "data": [{"employee": "Tristan Granier ", "link": "https://f....
```

Soon more command added please wait.

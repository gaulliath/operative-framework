package locales

var En_CA = map[string]interface{}{
	"internet": map[string]interface{}{
		"free_email": []string{
			"gmail.com", "yahoo.ca", "hotmail.com",
		},
		"domain_suffix": []string{
			"ca", "com", "biz", "info", "name", "net", "org",
		},
	},
	"phone_number": map[string]interface{}{
		"formats": []string{
			"###-###-####", "(###)###-####", "###.###.####", "1-###-###-####", "###-###-#### x###", "(###)###-#### x###", "1-###-###-#### x###", "###.###.#### x###", "###-###-#### x####", "(###)###-#### x####", "1-###-###-#### x####", "###.###.#### x####", "###-###-#### x#####", "(###)###-#### x#####", "1-###-###-#### x#####", "###.###.#### x#####",
		},
	},
	"address": map[string]interface{}{
		"postcode": "/[A-VX-Y][0-9][A-CEJ-NPR-TV-Z] ?[0-9][A-CEJ-NPR-TV-Z][0-9]/",
		"state": []string{
			"Alberta", "British Columbia", "Manitoba", "New Brunswick", "Newfoundland and Labrador", "Nova Scotia", "Northwest Territories", "Nunavut", "Ontario", "Prince Edward Island", "Quebec", "Saskatchewan", "Yukon",
		},
		"state_abbr": []string{
			"AB", "BC", "MB", "NB", "NL", "NS", "NU", "NT", "ON", "PE", "QC", "SK", "YK",
		},
		"default_country": []string{
			"Canada"}}}

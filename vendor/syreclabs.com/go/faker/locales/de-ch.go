package locales

var De_CH = map[string]interface{}{
	"address": map[string]interface{}{
		"country_code": []string{
			"CH", "CH", "CH", "DE", "AT", "US", "LI", "US", "HK", "VN",
		},
		"postcode": []string{
			"1###", "2###", "3###", "4###", "5###", "6###", "7###", "8###", "9###",
		},
		"default_country": []string{
			"Schweiz",
		},
	},
	"company": map[string]interface{}{
		"suffix": []string{
			"AG", "GmbH", "und Söhne", "und Partner", "& Co.", "Gruppe", "LLC", "Inc.",
		},
		"name": []string{
			"#{name.last_name} #{name.suffix}", "#{name.last_name}-#{name.last_name}", "#{name.last_name}, #{name.last_name} und #{name.last_name}",
		},
	},
	"internet": map[string]interface{}{
		"domain_suffix": []string{
			"com", "net", "biz", "ch", "de", "li", "at", "ch", "ch",
		},
	},
	"phone_number": map[string]interface{}{
		"formats": []string{
			"0800 ### ###", "0800 ## ## ##", "0## ### ## ##", "0## ### ## ##", "+41 ## ### ## ##", "0900 ### ###", "076 ### ## ##", "+4178 ### ## ##", "0041 79 ### ## ##"}}}

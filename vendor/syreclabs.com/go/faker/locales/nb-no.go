package locales

var Nb_NO = map[string]interface{}{
	"address": map[string]interface{}{
		"secondary_address": []string{
			"Leil. ###", "Oppgang A", "Oppgang B",
		},
		"postcode": []string{
			"####", "####", "####", "0###",
		},
		"street_address": []string{
			"#{address.street_name} #{address.building_number}",
		},
		"street_prefix": []string{
			"Øvre", "Nedre", "Søndre", "Gamle", "Østre", "Vestre",
		},
		"street_root": []string{
			"Eike", "Bjørke", "Gran", "Vass", "Furu", "Litj", "Lille", "Høy", "Fosse", "Elve", "Ku", "Konvall", "Soldugg", "Hestemyr", "Granitt", "Hegge", "Rogne", "Fiol", "Sol", "Ting", "Malm", "Klokker", "Preste", "Dam", "Geiterygg", "Bekke", "Berg", "Kirke", "Kors", "Bru", "Blåveis", "Torg", "Sjø",
		},
		"street_suffix": []string{
			"alléen", "bakken", "berget", "bråten", "eggen", "engen", "ekra", "faret", "flata", "gata", "gjerdet", "grenda", "gropa", "hagen", "haugen", "havna", "holtet", "høgda", "jordet", "kollen", "kroken", "lia", "lunden", "lyngen", "løkka", "marka", "moen", "myra", "plassen", "ringen", "roa", "røa", "skogen", "skrenten", "spranget", "stien", "stranda", "stubben", "stykket", "svingen", "tjernet", "toppen", "tunet", "vollen", "vika", "åsen",
		},
		"city": []string{
			"#{city_root}#{address.city_suffix}",
		},
		"common_street_suffix": []string{
			"sgate", "svei", "s Gate", "s Vei", "gata", "veien",
		},
		"building_number": []string{
			"#", "##",
		},
		"state": []string{
			"",
		},
		"default_country": []string{
			"Norge",
		},
		"city_root": []string{
			"Fet", "Gjes", "Høy", "Inn", "Fager", "Lille", "Lo", "Mal", "Nord", "Nær", "Sand", "Sme", "Stav", "Stor", "Tand", "Ut", "Vest",
		},
		"city_suffix": []string{
			"berg", "borg", "by", "bø", "dal", "eid", "fjell", "fjord", "foss", "grunn", "hamn", "havn", "helle", "mark", "nes", "odden", "sand", "sjøen", "stad", "strand", "strøm", "sund", "vik", "vær", "våg", "ø", "øy", "ås",
		},
		"street_name": []string{
			"#{street_root}#{address.street_suffix}", "#{street_prefix} #{street_root}#{address.street_suffix}", "#{name.first_name}#{common_street_suffix}", "#{name.last_name}#{common_street_suffix}",
		},
	},
	"company": map[string]interface{}{
		"suffix": []string{
			"Gruppen", "AS", "ASA", "BA", "RFH", "og Sønner",
		},
		"name": []string{
			"#{name.last_name} #{name.suffix}", "#{name.last_name}-#{name.last_name}", "#{name.last_name}, #{name.last_name} og #{name.last_name}",
		},
	},
	"internet": map[string]interface{}{
		"domain_suffix": []string{
			"no", "com", "net", "org",
		},
	},
	"name": map[string]interface{}{
		"prefix": []string{
			"Dr.", "Prof.",
		},
		"suffix": []string{
			"Jr.", "Sr.", "I", "II", "III", "IV", "V",
		},
		"name": []string{
			"#{name.prefix} #{name.first_name} #{name.last_name}", "#{name.first_name} #{name.last_name} #{name.suffix}", "#{feminine_name} #{feminine_name} #{name.last_name}", "#{masculine_name} #{masculine_name} #{name.last_name}", "#{name.first_name} #{name.last_name} #{name.last_name}", "#{name.first_name} #{name.last_name}",
		},
		"first_name": []string{
			"Emma", "Sara", "Thea", "Ida", "Julie", "Nora", "Emilie", "Ingrid", "Hanna", "Maria", "Sofie", "Anna", "Malin", "Amalie", "Vilde", "Frida", "Andrea", "Tuva", "Victoria", "Mia", "Karoline", "Mathilde", "Martine", "Linnea", "Marte", "Hedda", "Marie", "Helene", "Silje", "Leah", "Maja", "Elise", "Oda", "Kristine", "Aurora", "Kaja", "Camilla", "Mari", "Maren", "Mina", "Selma", "Jenny", "Celine", "Eline", "Sunniva", "Natalie", "Tiril", "Synne", "Sandra", "Madeleine", "Markus", "Mathias", "Kristian", "Jonas", "Andreas", "Alexander", "Martin", "Sander", "Daniel", "Magnus", "Henrik", "Tobias", "Kristoffer", "Emil", "Adrian", "Sebastian", "Marius", "Elias", "Fredrik", "Thomas", "Sondre", "Benjamin", "Jakob", "Oliver", "Lucas", "Oskar", "Nikolai", "Filip", "Mats", "William", "Erik", "Simen", "Ole", "Eirik", "Isak", "Kasper", "Noah", "Lars", "Joakim", "Johannes", "Håkon", "Sindre", "Jørgen", "Herman", "Anders", "Jonathan", "Even", "Theodor", "Mikkel", "Aksel",
		},
		"feminine_name": []string{
			"Emma", "Sara", "Thea", "Ida", "Julie", "Nora", "Emilie", "Ingrid", "Hanna", "Maria", "Sofie", "Anna", "Malin", "Amalie", "Vilde", "Frida", "Andrea", "Tuva", "Victoria", "Mia", "Karoline", "Mathilde", "Martine", "Linnea", "Marte", "Hedda", "Marie", "Helene", "Silje", "Leah", "Maja", "Elise", "Oda", "Kristine", "Aurora", "Kaja", "Camilla", "Mari", "Maren", "Mina", "Selma", "Jenny", "Celine", "Eline", "Sunniva", "Natalie", "Tiril", "Synne", "Sandra", "Madeleine",
		},
		"masculine_name": []string{
			"Markus", "Mathias", "Kristian", "Jonas", "Andreas", "Alexander", "Martin", "Sander", "Daniel", "Magnus", "Henrik", "Tobias", "Kristoffer", "Emil", "Adrian", "Sebastian", "Marius", "Elias", "Fredrik", "Thomas", "Sondre", "Benjamin", "Jakob", "Oliver", "Lucas", "Oskar", "Nikolai", "Filip", "Mats", "William", "Erik", "Simen", "Ole", "Eirik", "Isak", "Kasper", "Noah", "Lars", "Joakim", "Johannes", "Håkon", "Sindre", "Jørgen", "Herman", "Anders", "Jonathan", "Even", "Theodor", "Mikkel", "Aksel",
		},
		"last_name": []string{
			"Johansen", "Hansen", "Andersen", "Kristiansen", "Larsen", "Olsen", "Solberg", "Andresen", "Pedersen", "Nilsen", "Berg", "Halvorsen", "Karlsen", "Svendsen", "Jensen", "Haugen", "Martinsen", "Eriksen", "Sørensen", "Johnsen", "Myhrer", "Johannessen", "Nielsen", "Hagen", "Pettersen", "Bakke", "Skuterud", "Løken", "Gundersen", "Strand", "Jørgensen", "Kvarme", "Røed", "Sæther", "Stensrud", "Moe", "Kristoffersen", "Jakobsen", "Holm", "Aas", "Lie", "Moen", "Andreassen", "Vedvik", "Nguyen", "Jacobsen", "Torgersen", "Ruud", "Krogh", "Christiansen", "Bjerke", "Aalerud", "Borge", "Sørlie", "Berge", "Østli", "Ødegård", "Torp", "Henriksen", "Haukelidsæter", "Fjeld", "Danielsen", "Aasen", "Fredriksen", "Dahl", "Berntsen", "Arnesen", "Wold", "Thoresen", "Solheim", "Skoglund", "Bakken", "Amundsen", "Solli", "Smogeli", "Kristensen", "Glosli", "Fossum", "Evensen", "Eide", "Carlsen", "Østby", "Vegge", "Tangen", "Smedsrud", "Olstad", "Lunde", "Kleven", "Huseby", "Bjørnstad", "Ryan", "Rasmussen", "Nygård", "Nordskaug", "Nordby", "Mathisen", "Hopland", "Gran", "Finstad", "Edvardsen",
		},
	},
	"phone_number": map[string]interface{}{
		"formats": []string{
			"########", "## ## ## ##", "### ## ###", "+47 ## ## ## ##"}}}

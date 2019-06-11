package google_dork

import (
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"os"
)

type GoogleDorkModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushGoogleDorkModule(s *session.Session) *GoogleDorkModule{
	mod := GoogleDorkModule{
		sess: s,
		Stream: &s.Stream,
	}
	return &mod
}


func (module *GoogleDorkModule) Name() string{
	return "dorks_list"
}

func (module *GoogleDorkModule) Author() string{
	return "Tristan Granier"
}

func (module *GoogleDorkModule) Description() string{
	return "Lists available google dork"
}

func (module *GoogleDorkModule) GetType() string{
	return ""
}

func (module *GoogleDorkModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}


func (module *GoogleDorkModule) Start(){
	dorks := map[string]string{
		"cache:": "If you include other words in the query, Google will highlight those words within" +
			"the cached document. For instance, [cache:www.google.com web] will show the cached" +
			"content with the word “web” highlighted. This functionality is also accessible by" +
			"clicking on the “Cached” link on Google’s main results page. The query [cache:] will" +
			"show the version of the web page that Google has in its cache. For instance," +
			"[cache:www.google.com] will show Google’s cache of the Google homepage. Note there" +
			"can be no space between the “cache:” and the web page url.",
		"link:":"The query [link:] will list webpages that have links to the specified webpage." +
			"For instance, [link:www.google.com] will list webpages that have links pointing to the" +
			"Google homepage. Note there can be no space between the “link:” and the web page url.",
		"related:":"The query [related:] will list web pages that are “similar” to a specified web" +
			"page. For instance, [related:www.google.com] will list web pages that are similar to" +
			"the Google homepage. Note there can be no space between the “related:” and the web" +
			"page url.",
		"info:":"The query [info:] will present some information that Google has about that web" +
			"page. For instance, [info:www.google.com] will show information about the Google" +
			"homepage. Note there can be no space between the “info:” and the web page url.",
		"define:":"The query [define:] will provide a definition of the words you enter after it," +
			"gathered from various online sources. The definition will be for the entire phrase" +
			"entered (i.e., it will include all the words in the exact order you typed them).",
		"stocks:":"If you begin a query with the [stocks:] operator, Google will treat the rest" +
			"of the query terms as stock ticker symbols, and will link to a page showing stock" +
			"information for those symbols. For instance, [stocks: intc yhoo] will show information" +
			"about Intel and Yahoo. (Note you must type the ticker symbols, not the company name.)",
		"site:":"If you include [site:] in your query, Google will restrict the results to those" +
			"websites in the given domain. For instance, [help site:www.google.com] will find pages" +
			"about help within www.google.com. [help site:com] will find pages about help within" +
			".com urls. Note there can be no space between the “site:” and the domain.",
		"allintitle:":"If you start a query with [allintitle:], Google will restrict the results" +
			"to those with all of the query words in the title. For instance," +
			"allintitle: google search] will return only documents that have both “google" +
			"and “search” in the title.",
		"intitle:":"If you include [intitle:] in your query, Google will restrict the results" +
			"to documents containing that word in the title. For instance, [intitle:google search]" +
			"will return documents that mention the word “google” in their title, and mention the" +
			"word “search” anywhere in the document (title or no). Note there can be no space" +
			"between the “intitle:” and the following word. Putting [intitle:] in front of every" +
			"word in your query is equivalent to putting [allintitle:] at the front of your" +
			"query: [intitle:google intitle:search] is the same as [allintitle: google search].",
		"allinurl:":"If you start a query with [allinurl:], Google will restrict the results to" +
			"those with all of the query words in the url. For instance, [allinurl: google search]" +
			"will return only documents that have both “google” and “search” in the url. Note" +
			"that [allinurl:] works on words, not url components. In particular, it ignores" +
			"punctuation. Thus, [allinurl: foo/bar] will restrict the results to page with the" +
			"words “foo” and “bar” in the url, but won’t require that they be separated by a" +
			"slash within that url, that they be adjacent, or that they be in that particular" +
			"word order. There is currently no way to enforce these constraints.",
		"inurl:":"If you include [inurl:] in your query, Google will restrict the results to" +
			"documents containing that word in the url. For instance, [inurl:google search] will" +
			"return documents that mention the word “google” in their url, and mention the word" +
			"“search” anywhere in the document (url or no). Note there can be no space between" +
			"the “inurl:” and the following word. Putting “inurl:” in front of every word in your" +
			"query is equivalent to putting “allinurl:” at the front of your query:" +
			"[inurl:google inurl:search] is the same as [allinurl: google search].",
		"ext:":"If you include [ext:] in your query, Google will restrict the results" +
			"to documents containing this extension",
	}

	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{30,80})
	t.AppendHeader(table.Row{
		"Dork",
		"Resume",
	})
	for dork, resume := range dorks{
		t.AppendRow(table.Row{
			dork,
			resume,
		})
	}
	module.Stream.Render(t)
}
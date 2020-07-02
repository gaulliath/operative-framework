package linkedin_search

import (
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
)

type LinkedinSearchModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushLinkedinSearchModule(s *session.Session) *LinkedinSearchModule {
	mod := LinkedinSearchModule{
		sess:   s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Enterprise Name", "", true, session.STRING)
	mod.CreateNewParam("limit", "Limit search", "10", false, session.STRING)
	return &mod
}

func (module *LinkedinSearchModule) Name() string {
	return "linkedin.search"
}

func (module *LinkedinSearchModule) Description() string {
	return "Search employee on selected enterprise with Linkedin"
}

func (module *LinkedinSearchModule) Author() string {
	return "Tristan Granier"
}

func (module *LinkedinSearchModule) GetType() []string {
	return []string{
		session.T_TARGET_ENTERPRISE,
	}
}

func (module *LinkedinSearchModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *LinkedinSearchModule) Start() {
	opfClient := module.sess.Client
	paramEnterprise, _ := module.GetParameter("TARGET")
	target, err := module.sess.GetTarget(paramEnterprise.Value)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	paramLimit, _ := module.GetParameter("limit")
	url := "https://www.google.com/search?num=" + paramLimit.Value + "&start=0&hl=en&q=site:linkedin.com/in+" + strings.Replace(target.GetName(), " ", "+", -1)
	opfClient.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	res, err := opfClient.Perform("GET", url)
	if err != nil {
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		module.Stream.Error("A error as been occurred with a target.")
		return
	}

	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{0, 30})
	t.AppendHeader(table.Row{"Name", "Work", "Link"})

	resultFound := 0
	doc.Find(".g").Each(func(i int, s *goquery.Selection) {
		line := s.Find("h3").Text()
		line = strings.Replace(line, "| LinkedIn", "", 1)
		line = strings.Replace(line, "LinkedIn", "", 1)
		line = strings.Replace(line, "on LinkedIn", "", 1)
		if strings.Contains(line, "-") && len(strings.Split(strings.TrimSpace(line), "-")) > 1 {
			name := strings.Split(strings.TrimSpace(line), "-")[0]
			work := strings.Split(strings.TrimSpace(line), "-")[1]
			link, _ := s.Find("a[href]").First().Attr("href")
			t.AppendRow([]interface{}{name, work, link})
			result := target.NewResult()
			result.Set("Name", name)
			result.Set("Work", work)
			result.Set("Link", link)
			result.Save(module, target)
			resultFound = resultFound + 1

		}
	})
	module.Stream.Render(t)
	module.Stream.Success(strconv.Itoa(resultFound) + " employee(s) found with a Linkedin search.")
}

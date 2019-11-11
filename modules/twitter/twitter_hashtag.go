package twitter

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
)

type TwitterSearch struct {
	session.SessionModule
	Sess *session.Session `json:"-"`
}

func PushTwitterSearchModule(s *session.Session) *TwitterSearch {
	mod := TwitterSearch{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "Target search (hashtag, username, ...)", "", true, session.STRING)
	mod.CreateNewParam("RESULT_TYPE", "Specifies what type of search results you would prefer to receive. The current default is 'mixed'", "mixed", false, session.STRING)
	mod.CreateNewParam("SINCE", "SINCE DATE", time.Now().Format("2006-01-02 15:04:05"), false, session.STRING)
	mod.CreateNewParam("COUNT", "Number of tweets", "300", false, session.STRING)
	return &mod
}

func (module *TwitterSearch) Name() string {
	return "twitter.search"
}

func (module *TwitterSearch) Description() string {
	return "Get users who tweet the most from a hashtag"
}

func (module *TwitterSearch) Author() string {
	return "Tristan Granier"
}

func (module *TwitterSearch) GetType() string {
	return "twitter"
}

func (module *TwitterSearch) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *TwitterSearch) ReplaceOrAdd(target *session.Target, newResult session.TargetResults) {
	for moduleName, results := range target.Results {
		if moduleName == module.Name() {
			if len(results) > 0 {
				for k, result := range results {
					screenNameNew := strings.Split(newResult.Value, target.GetSeparator())[0]
					screenName := strings.Split(result.Value, target.GetSeparator())[0]
					if screenName == screenNameNew {
						target.Results[moduleName][k] = &newResult
						return
					}
				}
			}
		}
	}

	target.Save(module, newResult)
}

func (module *TwitterSearch) Start() {

	trg, err := module.GetParameter("TARGET")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	target, err := module.Sess.GetTarget(trg.Value)
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	resultType, err := module.GetParameter("RESULT_TYPE")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	maxCount, err := module.GetParameter("COUNT")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	v := url.Values{}
	v.Set("result_type", resultType.Value)
	v.Set("count", "100")

	api := anaconda.NewTwitterApiWithCredentials(module.Sess.Config.Twitter.Password, module.Sess.Config.Twitter.Api.SKey, module.Sess.Config.Twitter.Login, module.Sess.Config.Twitter.Api.Key)
	searchResult, err := api.GetSearch(target.Name, v)
	if err != nil {
		panic(err)
	}

	t := module.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{40, 0})
	t.AppendHeader(table.Row{
		"USERNAME",
		"NUMBERS",
	})

	search := make(map[string]int)
	tweets := make(map[string][]string)
	count := len(searchResult.Statuses)
	current := 0
	maxCountInt, _ := strconv.Atoi(maxCount.Value)
	for count >= 1 {
		for _, tweet := range searchResult.Statuses {
			if current >= maxCountInt {
				break
			}
			search[tweet.User.ScreenName] = search[tweet.User.ScreenName] + 1
			tweets[tweet.User.ScreenName] = append(tweets[tweet.User.ScreenName], tweet.Text)
			current = current + 1
		}
		searchResults, err := searchResult.GetNext(api)
		searchResult = searchResults
		if err != nil {
			panic(err)
		}
		time.Sleep(3 * time.Second)
		count = searchResult.Metadata.Count
		if current >= maxCountInt {
			break
		}
	}

	for username, value := range search {
		t.AppendRow(table.Row{
			username,
			value,
		})
		result := session.TargetResults{
			Header:    "screeName" + target.GetSeparator() + "count",
			Value:     username + target.GetSeparator() + strconv.Itoa(value),
			Auxiliary: tweets[username],
		}
		module.ReplaceOrAdd(target, result)
	}

	module.Sess.Stream.Render(t)
}

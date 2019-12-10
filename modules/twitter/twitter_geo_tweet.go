package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"net/url"
	"os"
)

type TwitterGeoTweet struct {
	session.SessionModule
	Sess *session.Session `json:"-"`
}

func PushTwitterGeoTweetModule(s *session.Session) *TwitterGeoTweet {
	mod := TwitterGeoTweet{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "TWITTER POST ID", "", true, session.STRING)
	return &mod
}

func (module *TwitterGeoTweet) Name() string {
	return "twitter.geo.search"
}

func (module *TwitterGeoTweet) Description() string {
	return "Show geolocation information of selected tweet"
}

func (module *TwitterGeoTweet) Author() string {
	return "Tristan Granier"
}

func (module *TwitterGeoTweet) GetType() string {
	return "text"
}

func (module *TwitterGeoTweet) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *TwitterGeoTweet) Start() {

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

	api := anaconda.NewTwitterApiWithCredentials(module.Sess.Config.Twitter.Password, module.Sess.Config.Twitter.Api.SKey, module.Sess.Config.Twitter.Login, module.Sess.Config.Twitter.Api.Key)
	v := url.Values{}
	var ids []int64
	ids = append(ids, module.Sess.StringToInt64(target.GetName()))
	search, err := api.GetTweetsLookupByIds(ids, v)
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	t := module.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	for _, tweet := range search {
		if tweet.Id == module.Sess.StringToInt64(target.GetName()) {
			latitude, _ := tweet.Latitude()
			longitude, _ := tweet.Longitude()
			t.AppendRow(table.Row{
				"ID",
				tweet.Id,
			})
			t.AppendRow(table.Row{
				"Text",
				tweet.Text,
			})
			t.AppendRow(table.Row{
				"Place Name",
				tweet.Place.FullName,
			})
			t.AppendRow(table.Row{
				"Place Type",
				tweet.Place.PlaceType,
			})
			t.AppendRow(table.Row{
				"Place Country",
				tweet.Place.Country,
			})
			t.AppendRow(table.Row{
				"Place URL",
				tweet.Place.URL,
			})
			t.AppendRow(table.Row{
				"Latitude",
				latitude,
			})
			t.AppendRow(table.Row{
				"Longitude",
				longitude,
			})
		}
	}

	module.Sess.Stream.Render(t)
}

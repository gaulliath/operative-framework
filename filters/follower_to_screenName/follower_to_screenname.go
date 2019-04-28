package follower_to_screenName

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"net/url"
	"os"
	"strconv"
)

type FollowerScreenName struct{
	session.SessionFilter
	Sess *session.Session
}

func PushFollowerScreenNameFilter(s *session.Session) *FollowerScreenName{
	mod := FollowerScreenName{
		Sess: s,
	}
	mod.AddModule("twitter_followers")
	mod.AddModule("twitter_following")
	return &mod
}

func (filter *FollowerScreenName) Name() string{
	return "follower_to_screen"
}

func (filter *FollowerScreenName) Description() string{
	return "Find screen name from twitter ID list"
}

func (filter *FollowerScreenName) Author() string{
	return "Tristan Granier"
}

func (filter *FollowerScreenName) Start(mod session.Module){
	api := anaconda.NewTwitterApiWithCredentials(filter.Sess.Config.Twitter.Password, filter.Sess.Config.Twitter.Api.SKey, filter.Sess.Config.Twitter.Login, filter.Sess.Config.Twitter.Api.Key)
	v := url.Values{}

	t := filter.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"id",
		"screen_name",
	})
	for _, id := range mod.GetResults(){
		id64, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			user, errU := api.GetUsersShowById(id64,v)
			if errU == nil{
				t.AppendRow(table.Row{
					id,
					user.ScreenName,

				})
			}
		}
	}
	filter.Sess.Stream.Render(t)
}

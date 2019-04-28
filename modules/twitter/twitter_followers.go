package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"net/url"
	"os"
	"strconv"
)
type TwitterFollower struct{
	session.SessionModule
	Sess *session.Session
}

func PushTwitterFollowerModule(s *session.Session) *TwitterFollower{
	mod := TwitterFollower{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "TWITTER USER SCREEN NAME", "", true, session.STRING)
	return &mod
}

func (module *TwitterFollower) Name() string{
	return "twitter_followers"
}

func (module *TwitterFollower) Description() string{
	return "Get followers from target user twitter account"
}

func (module *TwitterFollower) Author() string{
	return "Tristan Granier"
}

func (module *TwitterFollower) GetType() string{
	return "twitter"
}

func (module *TwitterFollower) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *TwitterFollower) Start(){

	var followerIds []int64

	trg, err := module.GetParameter("TARGET")
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}

	target, err := module.Sess.GetTarget(trg.Value)
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}


	api := anaconda.NewTwitterApiWithCredentials(module.Sess.Config.Twitter.Password, module.Sess.Config.Twitter.Api.SKey, module.Sess.Config.Twitter.Login, module.Sess.Config.Twitter.Api.Key)
	v := url.Values{}
	user, err := api.GetUserSearch(target.Name, v)
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}
	followers, err := api.GetFollowersUser(user[0].Id, v)
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}
	if followers.Next_cursor_str == "0"{
		for _, ids := range followers.Ids{
			followerIds = append(followerIds, ids)
		}
	}
	for followers.Next_cursor_str != "0"{
		for _, ids := range followers.Ids{
			followerIds = append(followerIds, ids)
		}
		v.Set("cursor", followers.Next_cursor_str)
		followers, err = api.GetFollowersUser(user[0].Id, v)
		if err != nil{
			module.Sess.Stream.Error(err.Error())
			break
		}
	}

	t := module.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Twitter ID",
	})
	separator := target.GetSeparator()
	for _, ids := range followerIds{
		t.AppendRow(table.Row{
			ids,
		})
		result := session.TargetResults{
			Header: "Twitter ID" + separator,
			Value: strconv.Itoa(int(ids)) + separator,
		}
		target.Save(module, result)
	}
	module.Sess.Stream.Render(t)
}
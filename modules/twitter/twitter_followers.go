package twitter

import (
	"fmt"
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
	mod.CreateNewParam("COUNT", "FOLLOWER LIMIT", "50", false, session.INT)
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


	argumentCount, errCount := module.GetParameter("COUNT")
	if errCount != nil{
		module.Sess.Stream.Error("Count parameters as not listed.")
		return
	}

	maxCount, errConv := strconv.Atoi(argumentCount.Value)
	if errConv != nil{
		module.Sess.Stream.Error("Error as occured with parameter 'COUNT'")
		return
	}
	current := 0
	if followers.Next_cursor_str == "0"{
		for _, ids := range followers.Ids{
			if current >= maxCount{
				break
			}
			followerIds = append(followerIds, ids)
			current = current + 1
		}
	}
	for followers.Next_cursor_str != "0"{
		for _, ids := range followers.Ids{
			fmt.Println(current)
			if current >= maxCount{
				break
			}
			followerIds = append(followerIds, ids)
			current = current + 1
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
	for _, ids := range followerIds{
		module.Results = append(module.Results, strconv.Itoa(int(ids)))
		t.AppendRow(table.Row{
			ids,
		})
		result := session.TargetResults{
			Header: "Twitter ID",
			Value: strconv.Itoa(int(ids)),
		}
		target.Save(module, result)
	}
	module.Sess.Stream.Render(t)
}
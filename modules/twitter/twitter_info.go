package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"net/url"
	"os"
	"strconv"
)

type TwitterInfo struct{
	session.SessionModule
	Sess *session.Session `json:"-"`
}

func PushTwitterInfoModule(s *session.Session) *TwitterInfo{
	mod := TwitterInfo{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "TWITTER USER SCREEN NAME", "", true, session.STRING)
	return &mod
}

func (module *TwitterInfo) Name() string{
	return "twitter.info"
}

func (module *TwitterInfo) Description() string{
	return "Get (re)tweets from target user twitter account"
}

func (module *TwitterInfo) Author() string{
	return "Tristan Granier"
}

func (module *TwitterInfo) GetType() string{
	return "twitter"
}

func (module *TwitterInfo) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *TwitterInfo) Start(){

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
	u := user[0]
	v.Set("screen_name", u.ScreenName)
	profile, err := api.GetUsersShow(target.GetName(), v)
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	t := module.Sess.Stream.GenerateTable()
	separator := target.GetSeparator()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{30, 30})
	t.AppendRow(table.Row{
		"USERNAME",
		profile.ScreenName,
	})
	t.AppendRow(table.Row{
		"FULLNAME",
		profile.Name,
	})
	t.AppendRow(table.Row{
		"PICS",
		profile.ProfileImageURL,
	})
	t.AppendRow(table.Row{
		"DESCRIPTION",
		profile.Description,
	})
	t.AppendRow(table.Row{
		"FOLLOWERS",
		profile.FollowersCount,
	})
	t.AppendRow(table.Row{
		"FOLLOWINGS",
		profile.FriendsCount,
	})
	t.AppendRow(table.Row{
		"EMAIL",
		profile.Email,
	})

	result := session.TargetResults{
		Header: "USERNAME" + separator + "FULLNAME" + separator + "PICS" + separator + "DESCRIPTION" + separator + "FOLLOWERS" + separator + "FOLLOWINGS" + separator + "EMAIL",
		Value:  profile.ScreenName + separator + profile.Name + separator + profile.ProfileImageURL + separator + profile.Description + separator + strconv.Itoa(profile.FollowersCount) + separator + strconv.Itoa(profile.FriendsCount) + separator + separator + profile.Email,
	}
	target.Save(module, result)

	module.Sess.Stream.Render(t)


}
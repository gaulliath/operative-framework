package instagram

import (
	"fmt"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"gopkg.in/ahmdrz/goinsta.v2"
	"os"
	"strconv"
)

type InstagramInfo struct{
	session.SessionModule
	Sess *session.Session
}

func PushInstagramInfoModule(s *session.Session) *InstagramInfo{
	mod := InstagramInfo{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "INSTAGRAM USER ACCOUNT", "",true,session.STRING)
	return &mod
}

func (module *InstagramInfo) Name() string{
	return "instagram_info"
}

func (module *InstagramInfo) Description() string{
	return "Get instagram account information"
}

func (module *InstagramInfo) Author() string{
	return "Tristan Granier"
}

func (module *InstagramInfo) GetType() string{
	return "instagram"
}

func (module *InstagramInfo) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *InstagramInfo) Start(){

	trg, err := module.GetParameter("TARGET")
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}

	target, err2 := module.Sess.GetTarget(trg.Value)
	if err2 != nil{
		fmt.Println(err2.Error())
		return
	}

	insta := goinsta.New(module.Sess.Config.Instagram.Login, module.Sess.Config.Instagram.Password)

	if err := insta.Login(); err != nil {
		fmt.Println(err)
		return
	}

	profil, err := insta.Profiles.ByName(target.Name)
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	t := module.Sess.Stream.GenerateTable()
	separator := target.GetSeparator()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{30, 30})
	t.AppendRow(table.Row{
		"USERNAME",
		profil.Username,
	})
	t.AppendRow(table.Row{
		"FULLNAME",
		profil.FullName,
	})
	t.AppendRow(table.Row{
		"PICS",
		profil.ProfilePicURL,
	})
	t.AppendRow(table.Row{
		"DESCRIPTION",
		profil.Biography,
	})
	t.AppendRow(table.Row{
		"FOLLOWERS",
		profil.FollowerCount,
	})
	t.AppendRow(table.Row{
		"FOLLOWINGS",
		profil.FollowingCount,
	})
	t.AppendRow(table.Row{
		"IS PRIVATE",
		profil.IsPrivate,
	})
	t.AppendRow(table.Row{
		"EMAIL",
		profil.Email,
	})

	result := session.TargetResults{
		Header: "USERNAME" + separator + "FULLNAME" + separator + "PICS" + separator + "DESCRIPTION" + separator + "FOLLOWERS" + separator + "FOLLOWINGS" + separator + "IS PRIVATE" + separator + "EMAIL",
		Value:  profil.Username + separator + profil.FullName + separator + profil.ProfilePicURL + separator + profil.Biography + separator + strconv.Itoa(profil.FollowerCount) + separator + strconv.Itoa(profil.FollowingCount) + separator + module.Sess.BooleanToString(profil.IsPrivate) + separator + profil.Email,
	}
	target.Save(module, result)

	module.Sess.Stream.Render(t)
}

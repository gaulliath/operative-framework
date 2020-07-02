package instagram

import (
	"fmt"
	"os"
	"strconv"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"gopkg.in/ahmdrz/goinsta.v2"
)

type InstagramInfo struct {
	session.SessionModule
	Sess *session.Session `json:"-"`
}

func PushInstagramInfoModule(s *session.Session) *InstagramInfo {
	mod := InstagramInfo{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "INSTAGRAM USER ACCOUNT", "", true, session.STRING)
	return &mod
}

func (module *InstagramInfo) Name() string {
	return "instagram.info"
}

func (module *InstagramInfo) Description() string {
	return "Get instagram account information"
}

func (module *InstagramInfo) Author() string {
	return "Tristan Granier"
}

func (module *InstagramInfo) GetType() []string {
	return []string{
		session.T_TARGET_INSTAGRAM,
	}
}

func (module *InstagramInfo) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *InstagramInfo) Start() {

	trg, err := module.GetParameter("TARGET")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	target, err2 := module.Sess.GetTarget(trg.Value)
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}

	insta := goinsta.New(module.Sess.Config.Instagram.Login, module.Sess.Config.Instagram.Password)

	if err := insta.Login(); err != nil {
		fmt.Println(err)
		return
	}

	profil, err := insta.Profiles.ByName(target.Name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t := module.Sess.Stream.GenerateTable()
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

	result := target.NewResult()
	result.Set("USERNAME", profil.Username)
	result.Set("FULLNAME", profil.FullName)
	result.Set("PICS", profil.ProfilePicURL)
	result.Set("DESCRIPTION", profil.Biography)
	result.Set("FOLLOWERS", strconv.Itoa(profil.FollowerCount))
	result.Set("FOLLOWING", strconv.Itoa(profil.FollowingCount))
	result.Set("IS PRIVATE", module.Sess.BooleanToString(profil.IsPrivate))
	result.Set("EMAIL", profil.Email)
	result.Save(module, target)

	module.Sess.Stream.Render(t)
}

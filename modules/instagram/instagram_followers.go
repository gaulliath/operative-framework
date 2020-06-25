package instagram

import (
	"fmt"
	"os"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"gopkg.in/ahmdrz/goinsta.v2"
)

type InstagramFollowers struct {
	session.SessionModule
	Sess *session.Session `json:"-"`
}

func PushInstagramFollowersModule(s *session.Session) *InstagramFollowers {
	mod := InstagramFollowers{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "INSTAGRAM USER ACCOUNT", "", true, session.STRING)
	return &mod
}

func (module *InstagramFollowers) Name() string {
	return "instagram.followers"
}

func (module *InstagramFollowers) Description() string {
	return "Get followers for target user instagram account"
}

func (module *InstagramFollowers) Author() string {
	return "Tristan Granier"
}

func (module *InstagramFollowers) GetType() []string {
	return []string{
		session.T_TARGET_INSTAGRAM,
	}
}

func (module *InstagramFollowers) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *InstagramFollowers) Start() {

	trg, err := module.GetParameter("TARGET")
	if err != nil {
		fmt.Println(err.Error())
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
	t.AppendHeader(table.Row{
		"Username",
		"Full Name",
	})
	followers := profil.Followers()
	for followers.Next() {
		for _, follower := range followers.Users {
			t.AppendRow(table.Row{
				follower.Username,
				follower.FullName,
			})
			result := target.NewResult()
			result.Set("username", follower.Username)
			result.Set("full_name", follower.FullName)
			result.Save(module, target)
		}
	}
	module.Sess.Stream.Render(t)
}

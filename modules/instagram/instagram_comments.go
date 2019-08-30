package instagram

import (
	"fmt"
	"os"
	"time"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"gopkg.in/ahmdrz/goinsta.v2"
)

type InstagramComments struct {
	session.SessionModule
	Sess *session.Session `json:"-"`
}

func PushInstagramCommentsModule(s *session.Session) *InstagramComments {
	mod := InstagramComments{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "INSTAGRAM USER ACCOUNT", "", true, session.STRING)
	return &mod
}

func (module *InstagramComments) Name() string {
	return "instagram.comments"
}

func (module *InstagramComments) Description() string {
	return "Get instagram comments from post."
}

func (module *InstagramComments) Author() string {
	return "Tristan Granier"
}

func (module *InstagramComments) GetType() string {
	return "instagram"
}

func (module *InstagramComments) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *InstagramComments) Start() {

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
	separator := target.GetSeparator()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{30, 30, 30, 30, 30})
	t.AppendHeader(table.Row{
		"ITEM ID",
		"USERNAME",
		"COMMENT",
	})
	media := profil.Feed()
	for media.Next() {
		if len(media.Items) > 0 {
			for _, item := range media.Items {
				item.Comments.Sync()
				for item.Comments.Next() {
					for _, c := range item.Comments.Items {
						t.AppendRow(table.Row{
							item.ID,
							c.User.Username,
							c.Text,
						})

						result := session.TargetResults{
							Header: "ITEM ID" + separator + "USERNAME" + separator + "COMMENT",
							Value:  item.ID + separator + c.User.Username + separator + c.Text,
						}

						target.Save(module, result)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}

	module.Sess.Stream.Render(t)
}

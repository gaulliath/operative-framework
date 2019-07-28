package instagram

import (
	"fmt"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"gopkg.in/ahmdrz/goinsta.v2"
	"os"
	"strconv"
	"time"
)

type InstagramFeed struct{
	session.SessionModule
	Sess *session.Session
}

func PushInstagramFeedModule(s *session.Session) *InstagramFeed{
	mod := InstagramFeed{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "INSTAGRAM USER ACCOUNT", "",true,session.STRING)
	mod.CreateNewParam("DOWNLOAD", "DOWNLOAD INSTAGRAM USER FEED", "false", false, session.BOOL)
	return &mod
}

func (module *InstagramFeed) Name() string{
	return "instagram_feed"
}

func (module *InstagramFeed) Description() string{
	return "Get instagram feed from selected target"
}

func (module *InstagramFeed) Author() string{
	return "Tristan Granier"
}

func (module *InstagramFeed) GetType() string{
	return "instagram"
}

func (module *InstagramFeed) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *InstagramFeed) Start(){

	trg, err := module.GetParameter("TARGET")
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}

	download, err := module.GetParameter("DOWNLOAD")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	hasDownload  := module.Sess.StringToBoolean(download.Value)
	exportPath := module.Sess.Config.Common.ExportDirectory + module.Name() + "/"
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
	t.SetAllowedColumnLengths([]int{30, 30, 30, 30, 30})
	t.AppendHeader(table.Row{
		"ID",
		"MEDIA",
		"CAPTION",
		"LIKES",
		"COMMENTS",
	})

	media := profil.Feed()
	downloaded := 0
	for media.Next() {
		if len(media.Items) > 0 {
			for _, item := range media.Items {
				t.AppendRow(table.Row{
					item.ID,
					item.Images.GetBest() + "...",
					item.Caption.Text,
					item.Likes,
					item.CommentCount,
				})
				result := session.TargetResults{
					Header: "ID" + separator + "MEDIA" + separator + "CAPTION" + separator + "LIKES" + separator + "COMMENTS",
					Value:  item.ID + separator + item.Images.GetBest() + separator + item.Caption.Text + separator + strconv.Itoa(item.Likes) + separator + strconv.Itoa(item.CommentCount),
				}
				target.Save(module, result)

				if hasDownload {
					_, _, err = item.Download(exportPath + target.Name, "")
					if err != nil {
						downloaded = downloaded + len(media.Items)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}
	module.Sess.Stream.Render(t)
	if hasDownload {
		module.Sess.Stream.Standard("Feed as exported at '" + exportPath + target.Name + "'")
	}
}

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

type InstagramFriends struct{
	session.SessionModule
	Sess *session.Session
	Friends map[string]string
}

func PushInstagramFriendsModule(s *session.Session) *InstagramFriends{
	mod := InstagramFriends{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "INSTAGRAM USER ACCOUNT", "",true,session.STRING)
	return &mod
}

func (module *InstagramFriends) Name() string{
	return "instagram_friends"
}

func (module *InstagramFriends) Description() string{
	return "Get possible friend on instagram (can take a time)"
}

func (module *InstagramFriends) Author() string{
	return "Tristan Granier"
}

func (module *InstagramFriends) GetType() string{
	return "instagram"
}

func (module *InstagramFriends) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *InstagramFriends) InFollowers(username string) bool{
	if _, ok := module.Friends[username]; ok {
		return true
	}
	return false
}

func (module *InstagramFriends) Start(){

	module.Friends = make(map[string]string)
	trg, err := module.GetParameter("TARGET")
	if err != nil{
		fmt.Println(err.Error())
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

	var friends []string

	followers := profil.Followers()
	for followers.Next(){
		for _, follower := range followers.Users{
			module.Friends[follower.Username] = follower.FullName
		}
		time.Sleep(1 * time.Second)
	}

	followings := profil.Following()
	for followings.Next(){
		for _, following := range followings.Users{
			if module.InFollowers(following.Username) {
				friends = append(friends, following.Username)
			}
		}
		time.Sleep(1 * time.Second)
	}

	t := module.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{30, 30, 30, 30, 30})
	t.AppendHeader(table.Row{
		"friends",
	})

	for _, username := range friends{
		t.AppendRow(table.Row{
			username,
		})

		result := session.TargetResults{
			Header: "FRIEND",
			Value:  username,
		}
		target.Save(module, result)
	}

	module.Sess.Stream.Render(t)
	module.Sess.Stream.Standard("Possible friend(s) '"+strconv.Itoa(len(friends))+"'")


}

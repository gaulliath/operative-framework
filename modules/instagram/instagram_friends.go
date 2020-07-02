package instagram

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"gopkg.in/ahmdrz/goinsta.v2"
)

type InstagramFriends struct {
	session.SessionModule
	Sess    *session.Session  `json:"-"`
	Friends map[string]string `json:"-"`
}

func PushInstagramFriendsModule(s *session.Session) *InstagramFriends {
	mod := InstagramFriends{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "INSTAGRAM USER ACCOUNT", "", true, session.STRING)
	mod.CreateNewParam("LIKES", "Search relationship including likes", "false", false, session.BOOL)
	mod.CreateNewParam("LIKE_ONLY", "Print only liked relationship", "false", false, session.BOOL)
	return &mod
}

func (module *InstagramFriends) Name() string {
	return "instagram.friends"
}

func (module *InstagramFriends) Description() string {
	return "Get possible friend on instagram (can take a time)"
}

func (module *InstagramFriends) Author() string {
	return "Tristan Granier"
}

func (module *InstagramFriends) GetType() []string {
	return []string{
		session.T_TARGET_INSTAGRAM,
	}
}

func (module *InstagramFriends) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *InstagramFriends) InFollowers(username string) bool {
	if _, ok := module.Friends[username]; ok {
		return true
	}
	return false
}

func (module *InstagramFriends) InSlice(slice []string, search string) bool {
	for _, element := range slice {
		if strings.ToLower(element) == strings.ToLower(search) {
			return true
		}
	}
	return false
}

func (module *InstagramFriends) Start() {

	module.Friends = make(map[string]string)

	var verifiedFriend []string

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

	likes, err := module.GetParameter("LIKES")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	onlyLike, err := module.GetParameter("LIKE_ONLY")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	withLike := module.Sess.StringToBoolean(likes.Value)
	oLike := module.Sess.StringToBoolean(onlyLike.Value)
	if oLike && !withLike {
		module.Sess.Stream.Error("You need 'LIKES' argument if you set 'LIKE_ONLY'")
		return
	}

	var friends []string

	followers := profil.Followers()
	for followers.Next() {
		for _, follower := range followers.Users {
			module.Friends[follower.Username] = follower.FullName
		}
		time.Sleep(1 * time.Second)
	}

	followings := profil.Following()
	for followings.Next() {
		for _, following := range followings.Users {
			if module.InFollowers(following.Username) {
				friends = append(friends, following.Username)
			}
		}
		time.Sleep(1 * time.Second)
	}

	if withLike {
		media := profil.Feed()
		for media.Next() {
			if len(media.Items) > 0 {
				for _, item := range media.Items {
					if item.Likes > 0 {
						for _, like := range item.Likers {
							if module.InSlice(friends, like.Username) {
								if !module.InSlice(verifiedFriend, like.Username) {
									verifiedFriend = append(verifiedFriend, like.Username)
								}
							}
						}
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
		if oLike {
			friends = verifiedFriend
		}
	}

	t := module.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{30, 2})
	if withLike {
		t.AppendHeader(table.Row{
			"friends",
			"<3",
		})
	} else {
		t.AppendHeader(table.Row{
			"friends",
		})
	}

	for _, username := range friends {
		if withLike {
			if module.InSlice(verifiedFriend, username) {
				t.AppendRow(table.Row{
					username,
					"X",
				})
			} else {
				t.AppendRow(table.Row{
					username,
					"-",
				})
			}
		} else {
			t.AppendRow(table.Row{
				username,
			})
		}

		result := target.NewResult()
		result.Set("FRIEND", username)
		result.Save(module, target)
	}

	module.Sess.Stream.Render(t)
	module.Sess.Stream.Standard("Possible friend(s) '" + strconv.Itoa(len(friends)) + "'")

}

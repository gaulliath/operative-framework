package phone_to_instagram

import (
	"fmt"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"gopkg.in/ahmdrz/goinsta.v2"
	"os"
)

type PhoneToInstagram struct{
	session.SessionFilter
	Sess *session.Session
}

func PushPhoneToInstagramFilter(s *session.Session) *PhoneToInstagram{
	mod := PhoneToInstagram{
		Sess: s,
	}
	mod.AddModule("phone_generator")
	return &mod
}

func (filter *PhoneToInstagram) Name() string{
	return "phone_to_instagram"
}

func (filter *PhoneToInstagram) Description() string{
	return "Find result of phone_generator module in instagram network."
}

func (filter *PhoneToInstagram) Author() string{
	return "Tristan Granier"
}

func (filter *PhoneToInstagram) Start(mod session.Module){
	insta := goinsta.New(filter.Sess.Config.Instagram.Login, filter.Sess.Config.Instagram.Password)

	if err := insta.Login(); err != nil {
		fmt.Println(err)
		return
	}

	var Contacts []goinsta.Contact
	Contacts = append(Contacts, goinsta.Contact{
		Numbers: mod.GetResults(),
	})
	syncAwser, err := insta.Contacts.SyncContacts(&Contacts)
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	t := filter.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{0,0,0,50})
	t.AppendHeader(table.Row{
		"Username",
		"isPrivate",
		"isVerified",
		"Picture",

	})
	if len(syncAwser.Users) > 0{
		for _,users := range syncAwser.Users{
			t.AppendRow(table.Row{
				users.Username,
				users.IsPrivate,
				users.IsVerified,
				users.ProfilePicURL,
			})
		}
	}

	filter.Sess.Stream.Render(t)

}

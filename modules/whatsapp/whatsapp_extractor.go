package whatsapp

import (
	"encoding/gob"
	"io/ioutil"
	"strconv"

	"encoding/json"
	"fmt"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"os"
	"time"
)

type ThumbUrl struct {
	EURL string `json:"eurl"`
    Tag string `json:"tag"`
    Status int64 `json:"status"`
}

type WhatsAppContacts struct{
	Contact whatsapp.Contact `json:"contact"`
	Picture ThumbUrl `json:"picture"`

}

type WhatsappExtractor struct{
	session.SessionModule
	Sess *session.Session
	Contacts []WhatsAppContacts `json:"contacts"`
}

func PushWhatsappExtractorModule(s *session.Session) *WhatsappExtractor{
	mod := WhatsappExtractor{
		Sess: s,
	}
	mod.CreateNewParam("FILE_PATH", "File to extract contacts", "contacts.json", true, session.STRING)
	return &mod
}

func (module *WhatsappExtractor) Name() string{
	return "whatsapp_extractor"
}

func (module *WhatsappExtractor) Description() string{
	return "Run reversed WhatsApp web & extract contacts"
}

func (module *WhatsappExtractor) Author() string{
	return "Tristan Granier"
}

func (module *WhatsappExtractor) GetType() string{
	return "whatsapp"
}

func (module *WhatsappExtractor) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func readSession() (whatsapp.Session, error) {
	s2 := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return s2, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&s2)
	if err != nil {
		return s2, err
	}
	return s2, nil
}

func (module *WhatsappExtractor) Start(){

	file_ouput, err := module.GetParameter("FILE_PATH")
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return

	}
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	//load saved session
	s2, err := readSession()
	if err == nil {
		//restore session
		s2, err = wac.RestoreWithSession(s2)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "restoring failed: %v\n", err)
			return
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
			module.Sess.Stream.Warning("Please scan this QRCode....")
		}()
		s2, err = wac.Login(qr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error during login: %v\n", err)
			return
		}
		module.Sess.Stream.Success("Successfully login, extracting contacts information please wait...")

		<-time.After(10 * time.Second)
		contacts := wac.Store.Contacts
		t := module.Sess.Stream.GenerateTable()
		t.SetOutputMirror(os.Stdout)
		t.SetAllowedColumnLengths([]int{0, 0, 30,})
		t.AppendHeader(table.Row{
			"Contact Name",
			"Contact JID",
			"Contact Picture",
		})

		max := 7000
		current := 1

		for _,v := range contacts{
			if current >= max{
				break
			}
			profilePicThumb, _ := wac.GetProfilePicThumb(v.Jid)
			profilePic := <- profilePicThumb
			Picture := ThumbUrl{}
			_ = json.Unmarshal([]byte(profilePic), &Picture)
			if Picture.EURL != "" {
				module.Contacts = append(module.Contacts, WhatsAppContacts{
					Contact: v,
					Picture: Picture,
				})
				t.AppendRow(table.Row{
					v.Name,
					v.Jid,
					Picture.EURL,
				})
				current = current + 1
			}
		}
		module.Sess.Stream.Render(t)
		module.Sess.Stream.Success("Total contacts : " + strconv.Itoa(len(contacts)))
		rankingsJson, _ := json.Marshal(module.Contacts)
		err = ioutil.WriteFile(file_ouput.Value, rankingsJson, 0644)
	}

}
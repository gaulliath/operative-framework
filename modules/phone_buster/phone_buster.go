package phone_buster

import (
	"fmt"
	"github.com/graniet/operative-framework/session"
	"github.com/segmentio/ksuid"
	"os"
	"strconv"
	"strings"
)

type PhoneBuster struct{
	session.SessionModule
	Sess *session.Session
	Current int
}

func PushPhoneBusterModule(s *session.Session) *PhoneBuster{
	mod := PhoneBuster{
		Sess: s,
		Current: 1,
	}

	mod.CreateNewParam("TARGET", "Target number without numbers to brute", "+1 (310) 999-33", true, session.STRING)
	mod.CreateNewParam("RANGE", "Max range to brute eg: 99, 999, 9999", "99", true, session.STRING)
	mod.CreateNewParam("FILE_PATH", "Location for generated VCards", "", false, session.STRING)
	return &mod
}

func (module *PhoneBuster) Name() string{
	return "phone_buster"
}

func (module *PhoneBuster) Description() string{
	return "Brute force phone number and generate VCard(s) (.vcf)"
}

func (module *PhoneBuster) Author() string{
	return "Tristan Granier"
}

func (module *PhoneBuster) GetType() string{
	return "phone"
}

func (module *PhoneBuster) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *PhoneBuster) Start(){
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

	argumentCount, err3 := module.GetParameter("RANGE")
	if err3 != nil{
		module.Sess.Stream.Error(err3.Error())
		return
	}

	maxRange, errConv := strconv.Atoi(argumentCount.Value)
	if errConv != nil{
		module.Sess.Stream.Error(errConv.Error())
		return
	}

	argumentFilePath, err2 :=  module.GetParameter("FILE_PATH")
	if err2 != nil{
		argumentFilePath = session.Param{
			Value: "",
		}
	}

	var file *os.File
	var errPath error

	if argumentFilePath.Value != ""{
		file, errPath = os.OpenFile(strings.TrimSpace(argumentFilePath.Value), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	} else {
		file, errPath = os.OpenFile("/Users/graniet/Desktop/VCARD/MRROBOT_1.vcf", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}
	if errPath != nil {
		fmt.Println(errPath.Error())
		return
	}

	minRange := 01

	for{

		number := target.Name + strconv.Itoa(minRange)
		module.Results = append(module.Results, number)
		module.Sess.Stream.Success(number)
		minRange = minRange + 1

		if minRange >= maxRange{
			break
		}
	}


	defer file.Close()
	for _, number := range module.Results{
		var uuid string
		uuid = "NY_" + ksuid.New().String()
		_, _ = file.WriteString("BEGIN:VCARD\nVERSION:3.0\nN:" + uuid + ";;;\nFN:" + uuid + "\nTEL;type=HOME:" + number + "\nEND:VCARD\n")
	}

	module.Sess.Stream.Success("Generation as been successfully ended")
	return
}
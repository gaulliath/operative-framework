package system

import (
	"github.com/graniet/operative-framework/session"
	"os"
	"os/exec"
)

type SystemModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushSystemModuleModule(s *session.Session) *SystemModule{
	mod := SystemModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("cmd", "Command to exec e.g: ls -l", "", true, session.STRING)
	return &mod
}

func (module *SystemModule) Name() string{
	return "sh"
}

func (module *SystemModule) Description() string{
	return "Execute system command"
}

func (module *SystemModule) Author() string{
	return "Tristan Granier"
}

func (module *SystemModule) GetType() string{
	return "command"
}

func (module *SystemModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *SystemModule) Start(){
	command, err := module.GetParameter("CMD")
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}

	cmd := exec.Command(command.Value)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}
	return
}

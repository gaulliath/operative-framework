package session_stream

import (
	"github.com/graniet/operative-framework/session"
)

type SessionStreamModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushSessionStreamModule(s *session.Session) *SessionStreamModule{
	mod := SessionStreamModule{
		sess:s,
		Stream: &s.Stream,
	}
	mod.CreateNewParam("VERBOSE","Change Verbosity (true/false)",mod.sess.BooleanToString(mod.Stream.Verbose),false,session.BOOL)
	mod.CreateNewParam("JSON", "Print a response with JSON format", mod.sess.BooleanToString(mod.Stream.JSON), false, session.BOOL)
	return &mod
}

func (module *SessionStreamModule) Name() string{
	return "session_stream"
}

func (module *SessionStreamModule) Description() string{
	return "Set a session event stream settings"
}

func (module *SessionStreamModule) Author() string{
	return "Tristan Granier"
}

func (module *SessionStreamModule) GetType() string{
	return "session"
}

func (module *SessionStreamModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *SessionStreamModule) Start(){
	paramVerbosity, err := module.GetParameter("VERBOSE")
	if err != nil{
		module.Stream.Error(err.Error())
		return
	}
	module.Stream.Verbose = module.sess.StringToBoolean(paramVerbosity.Value)
	return

}
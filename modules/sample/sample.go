package sample

import (
	"github.com/graniet/operative-framework/session"
)

type SampleModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushSampleModuleModule(s *session.Session) *SampleModule{
	mod := SampleModule{
		sess: s,
		Stream: &s.Stream,
	}
	mod.CreateNewParam("TARGET", "TWITTER ID e.g: 4378543", "", true, session.STRING)
	return &mod
}

func (module *SampleModule) Name() string{
	return "sample_module"
}

func (module *SampleModule) Description() string{
	return "Print hello with twitter ID"
}

func (module *SampleModule) Author() string{
	return "Tristan Granier"
}

func (module *SampleModule) GetType() string{
	return "twitter"
}

func (module *SampleModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *SampleModule) Start(){
	trg, err := module.GetParameter("TARGET")
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}

	target, err := module.sess.GetTarget(trg.Value)
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}
	module.sess.Stream.Standard("Hello " + target.GetName() + " \\o/")
	return
}

package cookies

import (
	"github.com/graniet/operative-framework/session"
)

type GetCookiesModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushGetCookiesModule(s *session.Session) *GetCookiesModule {
	mod := GetCookiesModule{
		sess:   s,
		Stream: &s.Stream,
	}
	mod.CreateNewParam("TARGET", "e.g: 127.0.0.1", "", true, session.BOOL)
	return &mod
}

func (module *GetCookiesModule) Name() string {
	return "get.cookies"
}

func (module *GetCookiesModule) Description() string {
	return "Get cookie from http header"
}

func (module *GetCookiesModule) Author() string {
	return "Tristan Granier"
}

func (module *GetCookiesModule) GetType() []string {
	return []string{
		session.T_TARGET_URL,
		session.T_TARGET_HEADER,
	}
}

func (module *GetCookiesModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *GetCookiesModule) Start() {
	trg, err := module.GetParameter("TARGET")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	target, err := module.sess.GetTarget(trg.Value)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	client := module.sess.Client
	res, err := client.Perform("GET", target.GetName())
	if err != nil {
		module.sess.NewEvent(session.ERROR_MODULE, err.Error())
		return
	}

	cookies := res.Header.Get("Set-Cookie")
	if cookies != "" {
		result := target.NewResult()
		result.Set("cookie", cookies)
		result.Save(module, target)
	}
	return
}

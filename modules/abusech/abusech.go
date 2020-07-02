package abusech

import (
	"github.com/graniet/operative-framework/session"
	"path"
	"strings"
)

type AbuseChModule struct {
	session.SessionModule
	sess   *session.Session `json:"-"`
	Stream *session.Stream  `json:"-"`
}

func PushAbuseChModule(s *session.Session) *AbuseChModule {
	mod := AbuseChModule{
		sess:   s,
		Stream: &s.Stream,
	}
	mod.CreateNewParam("TARGET", "e.g: 127.0.0.1", "", true, session.BOOL)
	return &mod
}

func (module *AbuseChModule) Name() string {
	return "check.abuse_ch"
}

func (module *AbuseChModule) Description() string {
	return "Check if target as malicious classified"
}

func (module *AbuseChModule) Author() string {
	return "Tristan Granier"
}

func (module *AbuseChModule) GetType() []string {
	return []string{
		session.T_TARGET_IP_ADDRESS,
	}
}

func (module *AbuseChModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *AbuseChModule) Start() {
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

	lists := []string{
		"https://feodotracker.abuse.ch/downloads/ipblocklist.txt",
		"https://sslbl.abuse.ch/blacklist/sslipblacklist.csv",
		"https://ransomwaretracker.abuse.ch/downloads/RW_DOMBL.txt",
	}

	for _, url := range lists {
		client := module.sess.Client
		req, err := client.Perform("GET", url)
		if err != nil {
			module.sess.Stream.Error(err.Error())
			return
		}

		content, err := client.Read(req.Body)
		if err != nil {
			module.Stream.Error(err.Error())
			return
		}

		if strings.Contains(string(content), target.GetName()) {
			result := target.NewResult()
			result.Set("source", path.Base(url))
			result.Save(module, target)
		}
	}
	return
}

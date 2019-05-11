package account_checker

import (
	"encoding/json"
	"github.com/graniet/operative-framework/session"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type AccountCheckerModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

type GeneratedSites struct {
	License []string `json:"license"`
	Authors []string `json:"authors"`
	Sites   []struct {
		Name                   string   `json:"name"`
		CheckURI               string   `json:"check_uri"`
		AccountExistenceCode   string   `json:"account_existence_code"`
		AccountExistenceString string   `json:"account_existence_string"`
		AccountMissingString   string   `json:"account_missing_string"`
		AccountMissingCode     string   `json:"account_missing_code"`
		KnownAccounts          []string `json:"known_accounts"`
		Category               string   `json:"category"`
		Valid                  bool     `json:"valid"`
		PrettyURI              string   `json:"pretty_uri,omitempty"`
		Comments               []string `json:"comments,omitempty"`
		KnownMissingAccounts   []string `json:"known_missing_accounts,omitempty"`
		AllowedTypes           []string `json:"allowed_types,omitempty"`
	} `json:"sites"`
}

func PushAccountCheckerModule(s *session.Session) *AccountCheckerModule{
	mod := AccountCheckerModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Username target for checking", "", true, session.STRING)
	return &mod
}


func (module *AccountCheckerModule) Name() string{
	return "account_checker"
}

func (module *AccountCheckerModule) Author() string{
	return "Tristan Granier"
}

func (module *AccountCheckerModule) Description() string{
	return "Identify the existence of a given acount on various sites"
}

func (module *AccountCheckerModule) GetType() string{
	return "username"
}

func (module *AccountCheckerModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *AccountCheckerModule) LoadingSites(sites *GeneratedSites) bool{
	req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/WebBreacher/WhatsMyName/master/web_accounts_list.json", nil)
	if err != nil{
		return false
	}
	client := http.Client{}
	req.Header.Set("User-Agent",	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	res, err := client.Do(req)
	if err != nil{
		module.Stream.Error(err.Error())
	} else {
		if res.StatusCode == 200 {
			bodyBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return false
			}
			_ = json.Unmarshal(bodyBytes, sites)
			return true
		}
	}
	return false
}

func (module *AccountCheckerModule) Start(){

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

	var sites GeneratedSites
	load := module.LoadingSites(&sites)
	if load == false{
		module.sess.Stream.Error("can't load website listing.")
		return
	}

	if len(sites.Sites) > 0{
		client := http.Client{}

		for _, site := range sites.Sites{
			url := strings.Replace(site.CheckURI, "{account}", target.GetName(), -1)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil{
				continue
			}
			module.sess.Stream.Standard("checking: " + url)
			req.Header.Set("User-Agent",	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
			res, err := client.Do(req)
			if err != nil{
				module.Stream.Error(err.Error())
			} else{
				code, err := strconv.Atoi(site.AccountExistenceCode)
				if err != nil{
					continue
				}
				if res.StatusCode == code{
					bodyBytes, err := ioutil.ReadAll(res.Body)
					if err != nil {
						continue
					}
					if strings.Contains(string(bodyBytes), site.AccountExistenceString){
						module.sess.Stream.Success("Account found : " + url)
						result := session.TargetResults{
							Header:"URL" + target.GetSeparator() + "WEBSITE",
							Value: url + target.GetSeparator() + site.Name,
						}
						target.Save(module, result)
						module.Results = append(module.Results, url)
					}
				}
			}
		}
	}
	module.sess.Stream.Success("'" + strconv.Itoa(len(sites.Sites)) + "' sites checked.")
	return
}

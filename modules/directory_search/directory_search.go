package directory_search

import (
	"bufio"
	"github.com/graniet/operative-framework/session"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type DirectorySearchModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushModuleDirectorySearch(s *session.Session) *DirectorySearchModule{

	mod := DirectorySearchModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "website target", "", true, session.STRING)
	mod.CreateNewParam("WORDLIST", "WORDLIST e.g: operative-framwork-default/directory_search/lists.txt", "", true, session.STRING)
	return &mod
}

func (module *DirectorySearchModule) Name() string{
	return "directory_search"
}

func (module *DirectorySearchModule) Author() string{
	return "Tristan Granier"
}

func (module *DirectorySearchModule) Description() string{
	return "Try to bust hidden directory"
}

func (module *DirectorySearchModule) GetType() string{
	return "website"
}

func (module *DirectorySearchModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *DirectorySearchModule) Start(){
	trg, err := module.GetParameter("TARGET")
	if err != nil{
		module.Stream.Error(err.Error())
		return
	}
	target, err := module.sess.GetTarget(trg.Value)
	if err != nil{
		module.Stream.Error(err.Error())
		return
	}

	if strings.Contains(target.GetName(), "://"){
		expProto := strings.Split(target.GetName(), "://")
		proto := expProto[0]
		expURL := ""
		if strings.Contains(target.GetName(), "/"){
			expURL = strings.Split(expProto[1], "/")[0]
			target.Name = proto + "://" + expURL + "/"
		}
	} else{

		if strings.Contains(target.GetName(), "/"){
			expURL := strings.Split(target.GetName(), "/")[0]
			target.Name = "https://" + expURL + "/"
		}
	}


	wordList, err := module.GetParameter("WORDLIST")
	if err != nil{
		module.Stream.Error(err.Error())
		return
	}

	file, errPath := os.Open(wordList.Value)
	if errPath != nil{
		module.Stream.Error(errPath.Error())
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	client := http.Client{}
	for scanner.Scan(){
		line := strings.TrimSpace(scanner.Text())
		url := target.GetName() + line
		req, err := http.NewRequest("GET", url, nil)
		if err != nil{
			continue
		}
		req.Header.Set("User-Agent",	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
		res, err := client.Do(req)
		if err != nil{
			module.Stream.Error(err.Error())
		} else{
			if res.StatusCode == 200{
				result := session.TargetResults{
					Header: "URL" + target.GetSeparator() + "STATUS",
					Value: url + target.GetSeparator() + strconv.Itoa(res.StatusCode),
				}
				module.Results = append(module.Results, url)
				target.Save(module, result)
				module.sess.Stream.Success(url + " : " + strconv.Itoa(res.StatusCode))
			} else if res.StatusCode == 404{
				result := session.TargetResults{
					Header: "URL" + target.GetSeparator() + "STATUS",
					Value: url + target.GetSeparator() + strconv.Itoa(res.StatusCode),
				}
				module.Results = append(module.Results, url)
				target.Save(module, result)
				module.sess.Stream.Standard(url + " : " + strconv.Itoa(res.StatusCode))
			} else{
				result := session.TargetResults{
					Header: "URL" + target.GetSeparator() + "STATUS",
					Value: url + target.GetSeparator() + strconv.Itoa(res.StatusCode),
				}
				module.Results = append(module.Results, url)
				target.Save(module, result)
				module.sess.Stream.Warning(url + " : " + strconv.Itoa(res.StatusCode))
			}
		}
	}
	return
}

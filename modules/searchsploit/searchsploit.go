package searchsploit

import (
	"encoding/json"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"os"
	"os/exec"
)

type SearchSploitModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

type Exploit struct {
	SEARCH         string `json:"SEARCH"`
	DBPATHEXPLOIT  string `json:"DB_PATH_EXPLOIT"`
	RESULTSEXPLOIT []struct {
		Title string `json:"Title"`
		URL   string `json:"URL"`
	} `json:"RESULTS_EXPLOIT"`
	DBPATHSHELLCODE  string        `json:"DB_PATH_SHELLCODE"`
	RESULTSSHELLCODE []interface{} `json:"RESULTS_SHELLCODE"`
	DBPATHPAPER      string        `json:"DB_PATH_PAPER"`
	RESULTSPAPER     []struct {
		Title string `json:"Title"`
		URL   string `json:"URL"`
	} `json:"RESULTS_PAPER"`
}

func PushSearchSploitModule(s *session.Session) *SearchSploitModule{
	mod := SearchSploitModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.WithProgram("searchsploit")
	mod.CreateNewParam("TARGET", "Name of software", "", true, session.STRING)
	return &mod
}

func (module *SearchSploitModule) Name() string{
	return "search_exploit"
}

func (module *SearchSploitModule) Description() string{
	return "Search exploit for specific software"
}

func (module *SearchSploitModule) Author() string{
	return "Tristan Granier"
}

func (module *SearchSploitModule) GetType() string{
	return "software"
}

func (module *SearchSploitModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *SearchSploitModule) Start(){
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

	output, err := exec.Command("searchsploit", "-w", "-o", "-j", target.GetName()).CombinedOutput()
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	var exploits Exploit

	err = json.Unmarshal(output, &exploits)
	if err != nil{
		module.sess.Stream.Error("No result found.")
		return
	}

	if len(exploits.RESULTSEXPLOIT) < 1{
		module.sess.Stream.Error("No result found.")
		return
	}

	t := module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{50, 30})
	t.AppendHeader(table.Row{
		"TITLE",
		"URL",
	})
	for _, exploit := range exploits.RESULTSEXPLOIT{
		t.AppendRow(table.Row{
			exploit.Title,
			exploit.URL,
		})

		result := session.TargetResults{
			Header: "TITLE" + target.GetSeparator() + "URL",
			Value: exploit.Title + target.GetSeparator() + exploit.URL,
		}
		target.Save(module, result)

		module.Results = append(module.Results, exploit.Title)
	}
	module.sess.Stream.Render(t)
}

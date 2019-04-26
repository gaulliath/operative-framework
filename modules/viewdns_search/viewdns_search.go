package viewdns_search

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/go-pretty/table"
	"net/http"
	"os"
	"strings"
)

type VDNSearch struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushWSearchModule(sess *session.Session) *VDNSearch{
	mod := VDNSearch{
		sess: sess,
		Stream: &sess.Stream,
	}

	mod.CreateNewParam("TARGET", "Email Address","", true, session.STRING)
	return &mod
}

func (module *VDNSearch) Name() string{
	return "email_to_domain"
}

func (module *VDNSearch) Description() string{
	return "Find possible linked website with email address"
}

func (module *VDNSearch) Author() string{
	return "Tristan Granier"
}

func (module *VDNSearch) GetType() string{
	return "email"
}

func (module *VDNSearch) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *VDNSearch) Start(){
	targetId, err := module.GetParameter("TARGET")
	if err != nil{
		module.Stream.Error(err.Error())
		return
	}

	target, err := module.sess.GetTarget(targetId.Value)
	if err != nil{
		module.Stream.Error(err.Error())
		return
	}
	if target.GetType() != module.GetType(){
		module.Stream.Error("Target with type '"+target.GetType()+"' isn't valid module need '"+module.GetType()+"' type.")
		return
	}

	separator := target.GetSeparator()

	url := "https://viewdns.info/reversewhois/?q="+target.GetName()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent",	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		module.Stream.Error("Argument 'URL' can't be reached.")
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		module.Stream.Error(err.Error())
		return
	}

	headerReturn := "Domain Name" + separator + "Creation Date" + separator + "Registrar"
	resultFound := 0
	t := module.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Domain Name",
		"Creation Date",
		"Registrar",
	})
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("border")
		if exist && val == "1"{
			s.Find("tr").Each(func(i int, tr *goquery.Selection){
				valueReturn := ""
				resRow := table.Row{}
				tr.Find("td").Each(func(i int, td *goquery.Selection){
					element := strings.TrimSpace(td.Text())
					if element != "Domain Name" && element != "Creation Date" && element != "Registrar"{
						if valueReturn == ""{
							valueReturn = valueReturn + element
						} else{
							valueReturn = valueReturn + separator + element
						}
						resRow = append(resRow, element)
					}
				})
				if valueReturn != ""{
					result := session.TargetResults{
						Header: headerReturn,
						Value: valueReturn,
					}
					target.Save(module, result)
					t.AppendRow(resRow)
					resultFound = resultFound + 1
				}
			})
		}
	})
	module.Stream.Render(t)

}
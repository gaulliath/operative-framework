package report

import (
	"encoding/json"
	"github.com/graniet/operative-framework/session"
	"io/ioutil"
)

type ReportJSON struct{
	session.SessionModule
	Sess *session.Session `json:"-"`
}

func PushReportJSONModule(s *session.Session) *ReportJSON{
	mod := ReportJSON{
		Sess: s,
	}

	mod.CreateNewParam("EXPORT_FILE", "Export file e.g: /path/to/json", "./report.json", true, session.STRING)
	return &mod
}

func (module *ReportJSON) Name() string{
	return "report.json"
}

func (module *ReportJSON) Description() string{
	return "Generate session report to JSON format"
}

func (module *ReportJSON) Author() string{
	return "Tristan Granier"
}

func (module *ReportJSON) GetType() string{
	return ""
}

func (module *ReportJSON) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *ReportJSON) Start(){
	reportOutput, err := module.GetParameter("EXPORT_FILE")
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}

	file, _ := json.MarshalIndent(module.Sess.ExportNow(), "", " ")

	_ = ioutil.WriteFile(reportOutput.Value, file, 0755)

	module.Sess.Stream.Success("Report as generated to '" + reportOutput.Value + "'")
	return
}

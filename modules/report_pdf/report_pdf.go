package report_pdf

import (
	"fmt"
	"github.com/graniet/operative-framework/session"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
	"strconv"
	"strings"
)

type ReportPDF struct{
	session.SessionModule
	Sess *session.Session
}

func PushReportPDFModule(s *session.Session) *ReportPDF{
	mod := ReportPDF{
		Sess: s,
	}

	mod.CreateNewParam("EXPORT_FILE", "Export file e.g: /path/to/pdf", "./report.pdf", false, session.STRING)
	return &mod
}

func (module *ReportPDF) Name() string{
	return "report_pdf"
}

func (module *ReportPDF) Description() string{
	return "Generate session report to PDF format"
}

func (module *ReportPDF) Author() string{
	return "Tristan Granier"
}

func (module *ReportPDF) GetType() string{
	return ""
}

func (module *ReportPDF) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *ReportPDF) Start(){
	reportOutput, err := module.GetParameter("EXPORT_FILE")
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	url := "http://tristan-granier.com/static/img/operative.png"
	httpimg.Register(pdf, url, "")
	pdf.Image(url, 10, 5, 90, 20, false,"", 0, "")

	TextSeparator := func(text string, size float64){
		// Arial 12
		pdf.SetFont("Arial", "", size)
		// Background color
		pdf.SetFillColor(230, 126, 0)
		pdf.SetTextColor(255, 255, 255)
		pdf.Ln(20)
		// Title
		pdf.CellFormat(0, 6, fmt.Sprintf(text),
			"", 1, "L", true, 0, "")
		// Line break
		pdf.Ln(4)
	}

	TextSeparator("Session: " + module.Sess.SessionName, 12)
	Reset := func(){
		pdf.SetFont("Arial", "", 12)
		pdf.SetFillColor(255, 255, 255)
		pdf.SetTextColor(0, 0, 0)
	}
	Reset()
	pdf.CellFormat(0, 6,fmt.Sprintf("This is a report generated on 10-01-2019 with the operative framework."), "",1,"L",true,0,"")
	pdf.CellFormat(0,6, fmt.Sprintf("All data below is confidential"), "", 1, "L", true,0,"")
	Title := func(text string){
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 18)
		pdf.CellFormat(0, 6, fmt.Sprintf(text),"", 1, "L",true,0, "")
		pdf.SetFont("Arial", "", 12)
		pdf.Ln(1)
	}
	SubTitle := func(text string){
		pdf.CellFormat(0,10, fmt.Sprintf(text), "", 1, "L", true,0,"")
		pdf.Ln(8)
	}
	Targets := func() {
		pdf.CellFormat(80, 7, "TARGET", "1", 0, "", false, 0, "")
		pdf.CellFormat(80, 7, "TYPE", "1", 0, "", false, 0, "")
		pdf.Ln(-1)
		for _, target := range module.Sess.Targets {
			pdf.CellFormat(80, 6, target.Name, "1", 0, "", false, 0, "")
			pdf.CellFormat(80, 6, target.Type, "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}
	LinkedTargets := func(target *session.Target){
		pdf.CellFormat(80, 7, "TARGET", "1", 0, "", false, 0, "")
		pdf.CellFormat(80, 7, "NAME", "1", 0, "", false, 0, "")
		pdf.Ln(-1)

		if len(target.TargetLinked) > 0 {
			for _, link := range target.TargetLinked {
				trg, err := module.Sess.GetTarget(link.TargetId)
				if err == nil {
					pdf.CellFormat(80, 6, trg.GetId(), "1", 0, "", false, 0, "")
					pdf.CellFormat(80, 6, trg.GetName(), "1", 0, "", false, 0, "")
					pdf.Ln(-1)
				}
			}
		} else{
			pdf.CellFormat(80, 6, "No target(s) linked", "1", 0, "", false, 0, "")
			pdf.CellFormat(80, 6, "", "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}

	}
	Modules := func() {

		pdf.CellFormat(50, 7, "NAME", "1", 0, "", false, 0, "")
		pdf.CellFormat(120, 7, "DESCRIPTION", "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 7, "TYPE", "1", 0, "", false, 0, "")
		pdf.Ln(-1)
		for _, mod := range module.Sess.Modules {
			pdf.CellFormat(50, 6, mod.Name(), "1", 0, "", false, 0, "")
			pdf.CellFormat(120, 6, mod.Description(), "1", 0, "", false, 0, "")
			pdf.CellFormat(25, 6, mod.GetType(), "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}

	ViewResults := func(target *session.Target){
		pdf.CellFormat(80, 7, "NAME", "1", 0, "", false, 0, "")
		pdf.CellFormat(80, 7, "RESULT(S)", "1", 0, "", false, 0, "")
		pdf.Ln(-1)
		for mod, result := range target.Results {
			pdf.CellFormat(80, 6, mod, "1", 0, "", false, 0, "")
			pdf.CellFormat(80, 6, strconv.Itoa(len(result)), "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}

	PrintNote := func(notes []session.Note){
		for _, note := range notes {
			pdf.SetFillColor(230, 126, 0)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFont("Arial", "", 12)
			pdf.CellFormat(0, 6, note.Text, "1", 1, "L", true, 0, "")
			pdf.Ln(2)
			Reset()
		}
	}

	PrintResults := func(target *session.Target){
		for mod, result := range target.Results {
			Title(mod)
			SubTitle("results listed bellow:")
			if len(result) > 0{
				for _, res := range result{
					headers := strings.Split(result[0].Header, target.GetSeparator())
					values := strings.Split(res.Value, target.GetSeparator())
					for k, h := range headers{
						pdf.CellFormat(80, 7, h, "1", 0, "", false, 0, "")
						pdf.CellFormat(110, 7, values[k], "1", 0, "", false, 0, "")
						pdf.Ln(-1)

					}

					if len(res.Notes) > 0{
						PrintNote(res.Notes)
					}
					pdf.Ln(-1)
				}
			}
		}
	}

	Title("MODULE(S)")
	SubTitle("Module(s) loaded with session")
	Modules()

	pdf.AddPage()

	Title("TARGET(S)")
	SubTitle("The targets of the session are listed below")
	Targets()

	for _, target := range module.Sess.Targets{
		TextSeparator("TARGET: " + target.GetId() + " - " + target.GetName(), 14)
		Reset()
		Title("LINKED TARGET(S)")
		SubTitle("You can see if below the targets link automatically by operative.")

		LinkedTargets(target)

		Title("MODULE(S) EXECUTED")
		SubTitle("List of modules executed with this targets")
		ViewResults(target)
		PrintResults(target)

	}





	err = pdf.OutputFileAndClose(reportOutput.Value)
	if err != nil{
		module.Sess.Stream.Error(err.Error())
		return
	}
	module.Sess.Stream.Success("Report as generated to '" + reportOutput.Value + "'")
	return
}

package info_greffe

import (
	"encoding/json"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"time"
)

type InfoGreffeRegistration struct{
	session.SessionModule
	sess *session.Session `json:"-"`
	Stream *session.Stream `json:"-"`
}

type InfoGreffe struct {
	Nhits      int `json:"nhits"`
	Parameters struct {
		Dataset  []string `json:"dataset"`
		Timezone string   `json:"timezone"`
		Rows     int      `json:"rows"`
		Format   string   `json:"format"`
		Facet    []string `json:"facet"`
	} `json:"parameters"`
	Records []struct {
		Datasetid string `json:"datasetid"`
		Recordid  string `json:"recordid"`
		Fields    struct {
			Departement         string    `json:"departement"`
			Ville               string    `json:"ville"`
			Siren               string    `json:"siren"`
			CodePostal          string    `json:"code_postal"`
			DateDePublication   string    `json:"date_de_publication"`
			Statut              string    `json:"statut"`
			Nic                 string    `json:"nic"`
			CodeApe             string    `json:"code_ape"`
			Adresse             string    `json:"adresse"`
			NumDept             string    `json:"num_dept"`
			LibelleApe          string    `json:"libelle_ape"`
			Greffe              string    `json:"greffe"`
			DateImmatriculation string    `json:"date_immatriculation"`
			FormeJuridique      string    `json:"forme_juridique"`
			Geolocalisation     []float64 `json:"geolocalisation"`
			Denomination        string    `json:"denomination"`
			Region              string    `json:"region"`
			FicheIdentite       string    `json:"fiche_identite"`
		} `json:"fields,omitempty"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		RecordTimestamp time.Time `json:"record_timestamp"`
	} `json:"records"`
	FacetGroups []struct {
		Name   string `json:"name"`
		Facets []struct {
			Name  string `json:"name"`
			Path  string `json:"path"`
			Count int    `json:"count"`
			State string `json:"state"`
		} `json:"facets"`
	} `json:"facet_groups"`
}

func PushInfoGreffeRegistrationModule(s *session.Session) *InfoGreffeRegistration{
	mod := InfoGreffeRegistration{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "keyword", "", true, session.STRING)
	return &mod
}

func (module *InfoGreffeRegistration) Name() string{
	return "info_greffe.registration"
}

func (module *InfoGreffeRegistration) Description() string{
	return "Search enterprise registration in info greffe open data"
}

func (module *InfoGreffeRegistration) Author() string{
	return "Tristan Granier"
}

func (module *InfoGreffeRegistration) GetType() string{
	return "text"
}


func (module *InfoGreffeRegistration) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *InfoGreffeRegistration) Start(){
	keywordTarget, err := module.GetParameter("TARGET")
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	keyword, err := module.sess.GetTarget(keywordTarget.Value)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	url := "https://opendata.datainfogreffe.fr/api/records/1.0/search/?dataset=entreprises-immatriculees-2015&q="+url2.QueryEscape(keyword.Name)+"&facet=libelle&facet=forme_juridique&facet=code_postal&facet=ville&facet=region&facet=greffe&facet=date"
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	var results InfoGreffe
	err = json.Unmarshal(body, &results)
	if err != nil {
		module.Stream.Error(err.Error())
		return
	}

	separator := keyword.GetSeparator()

	for _, record := range results.Records {
		t := module.Stream.GenerateTable()
		t.SetOutputMirror(os.Stdout)
		t.SetAllowedColumnLengths([]int{0, 30,})

		t.AppendRow(table.Row{
			"NAME",
			record.Fields.Denomination,
		})
		t.AppendRow(table.Row{
			"SIREN",
			record.Fields.Siren,
		})
		t.AppendRow(table.Row{
			"ADDRESS",
			record.Fields.Adresse,
		})
		t.AppendRow(table.Row{
			"CITY",
			record.Fields.Ville,
		})
		t.AppendRow(table.Row{
			"CITY CODE",
			record.Fields.CodePostal,
		})
		t.AppendRow(table.Row{
			"REGION",
			record.Fields.Region,
		})
		if len(record.Fields.Geolocalisation) > 1 {
			t.AppendRow(table.Row{
				"LAT",
				record.Fields.Geolocalisation[0],
			})
			t.AppendRow(table.Row{
				"LNG",
				record.Fields.Geolocalisation[1],
			})
		}
		t.AppendRow(table.Row{
			"GREFFE",
			record.Fields.Greffe,
		})
		t.AppendRow(table.Row{
			"STATUS",
			record.Fields.Statut,
		})

		module.Stream.Render(t)

		result := session.TargetResults{
			Header: "NAME" + separator + "SIREN" + separator + "ADDRESS" + separator + "CITY",
			Value: record.Fields.Denomination + separator + record.Fields.Siren + separator + record.Fields.Adresse + separator + record.Fields.Ville,
		}
		keyword.Save(module, result)
	}

	module.Stream.Success(strconv.Itoa(len(results.Records)) + " records found.")
}

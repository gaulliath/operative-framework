package ip_information

import (
	"encoding/json"
	"fmt"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type IpInformation struct {
	session.SessionModule
	Sess *session.Session `json:"-"`
}

type Information struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

func PushIpInformationModule(s *session.Session) *IpInformation {
	mod := IpInformation{
		Sess: s,
	}

	mod.CreateNewParam("TARGET", "IP Address", "", true, session.STRING)
	return &mod
}

func (module *IpInformation) Name() string {
	return "ip.info"
}

func (module *IpInformation) Description() string {
	return "Get information from IP Address"
}

func (module *IpInformation) Author() string {
	return "Tristan Granier"
}

func (module *IpInformation) GetType() []string {
	return []string{
		session.T_TARGET_IP_ADDRESS,
	}
}

func (module *IpInformation) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *IpInformation) Start() {

	trg, err := module.GetParameter("TARGET")
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	target, err := module.Sess.GetTarget(trg.Value)
	if err != nil {
		module.Sess.Stream.Error(err.Error())
		return
	}

	url := "https://freegeoip.app/json/" + target.Name
	var IpAddressInformation Information

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &IpAddressInformation)
	if err != nil {
		module.Sess.Stream.Error(err.Error())
	}

	t := module.Sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{40, 0})
	t.AppendRow(table.Row{
		"IP",
		IpAddressInformation.IP,
	})
	t.AppendRow(table.Row{
		"CountryCode",
		IpAddressInformation.CountryCode,
	})
	t.AppendRow(table.Row{
		"CountryName",
		IpAddressInformation.CountryName,
	})
	t.AppendRow(table.Row{
		"RegionCode",
		IpAddressInformation.RegionCode,
	})
	t.AppendRow(table.Row{
		"RegionName",
		IpAddressInformation.RegionName,
	})
	t.AppendRow(table.Row{
		"City",
		IpAddressInformation.City,
	})
	t.AppendRow(table.Row{
		"ZipCode",
		IpAddressInformation.ZipCode,
	})
	t.AppendRow(table.Row{
		"TimeZone",
		IpAddressInformation.TimeZone,
	})
	t.AppendRow(table.Row{
		"Latitude",
		IpAddressInformation.Latitude,
	})
	t.AppendRow(table.Row{
		"Longitude",
		IpAddressInformation.Longitude,
	})
	t.AppendRow(table.Row{
		"MetroCode",
		IpAddressInformation.MetroCode,
	})

	module.Sess.Stream.Render(t)

	result := target.NewResult()
	result.Set("IP", IpAddressInformation.IP)
	result.Set("CountryCode", IpAddressInformation.CountryCode)
	result.Set("RegionCode", IpAddressInformation.RegionCode)
	result.Set("RegionName", IpAddressInformation.RegionName)
	result.Set("City", IpAddressInformation.City)
	result.Set("ZipCode", IpAddressInformation.ZipCode)
	result.Set("TimeZone", IpAddressInformation.TimeZone)
	result.Set("Latitude", fmt.Sprintf("%f", IpAddressInformation.Latitude))
	result.Set("Longitude", fmt.Sprintf("%f", IpAddressInformation.Longitude))
	result.Set("MetroCode", strconv.Itoa(IpAddressInformation.MetroCode))
	result.Save(module, target)
}

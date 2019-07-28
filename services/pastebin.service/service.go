//	AUTHOR	: 	TRISTAN GRANIER
//	RESUME	:	Search matching element to pastebin
//	TIME	:	Every 24 hours
package pastebin_service

import (
	"bytes"
	"encoding/json"
	"github.com/graniet/operative-framework/session"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type Service struct{
	session.Service
	session	*session.Session
}

func GetNewService(sess *session.Session) *Service{
	return &Service{
		session: sess,
	}
}

// Service as been started every 24 hours
func (service *Service) GetHibernate() time.Duration{
	return 24 * time.Hour
}

// Service name as been set here
func (service *Service) Name() string{
	return "pastebin.service"
}

// Define if service need configuration file
func (service *Service) HasConfiguration() bool{
	return true
}

// Get configuration
func (service *Service) GetConfiguration() map[string]string{
	configuration := make(map[string]string)
	configuration["MATCH"] = "example@gmail.com,example2@gmail.com"
	configuration["TO_SERVER"] = "false"
	configuration["SERVER_URI"] = "http://example.com/api/v1.0/insert"
	configuration["VERBOSE"] = "true"
	return configuration
}

// Get required fields in configuration file
func (service *Service) GetRequired() []string{
	return []string{
		"MATCH",
		"TO_SERVER",
		"VERBOSE",
	}
}

// Fetching matching to pastebin
func (service *Service) Fetch(configuration map[string]string, match string) (bool, error){
	module, err := service.session.SearchModule("pastebin")
	if err != nil {
		return false, err
	}
	targetId, err := service.session.AddTarget("text", match)
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") || targetId == ""{
			return false, err
		}
	}

	_, err = module.SetParameter("TARGET", targetId)
	if err != nil {
		return false, err
	}

	_, err = module.SetParameter("LIMIT", "10")
	if err != nil {
		return false, err
	}

	module.Start()

	target, err := service.session.GetTarget(targetId)
	if err != nil {
		return false, err
	}

	results, err := target.GetFormatedResults("pastebin")

	js, err := json.Marshal(&results)
	if err != nil{
		return false, err
	}

	if strings.ToLower(configuration["TO_SERVER"]) == "true" {
		if strings.ToLower(configuration["VERBOSE"]) == "true" {
			log.Println("Prepare request '"+match+"' at '"+configuration["SERVER_URI"]+"'")
		}
		req, err := http.NewRequest("POST", configuration["SERVER_URI"], bytes.NewBuffer(js))
		if err != nil {
			return false, errors.New("Can't make request to '"+configuration["SERVER_URI"]+"'")
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		if strings.ToLower(configuration["VERBOSE"]) == "true" {
			log.Println("Request '"+match+"' as been sent at '"+configuration["SERVER_URI"]+"'")
		}
	}
	return true, nil
}

// Service Execution with opf routine
func (service *Service) Run() (bool, error){

	configuration, _ := godotenv.Read(service.session.Config.Common.ConfigurationService + service.Name() + "/service.conf")
	service.session.Stream.Verbose = false

	if strings.Contains(configuration["MATCH"], ",") {
		matches := strings.Split(configuration["MATCH"], ",")
		for _, match := range matches{
			ret, err := service.Fetch(configuration, match)
			if err != nil {
				return ret, err
			}
		}
	} else {
		ret, err := service.Fetch(configuration, configuration["MATCH"])
		if err != nil {
			return ret, err
		}
	}
	return true, nil
}
//	AUTHOR	: 	TRISTAN GRANIER
//	RESUME	: 	This service get domain name registered with email address
//	TIME	:	Every 48 hours
package email_to_domain_service

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

// Service as been started every 48 hours
func (service *Service) GetHibernate() time.Duration{
	return 48 * time.Hour
}

// Service name as been set here
func (service *Service) Name() string{
	return "email_to_domain.service"
}

// Define if service need configuration file
func (service *Service) HasConfiguration() bool{
	return true
}

// Get configuration
func (service *Service) GetConfiguration() map[string]string{
	configuration := make(map[string]string)
	configuration["EMAIL"] = "example@email.com,example2@email.com"
	configuration["TO_SERVER"] = "false"
	configuration["SERVER_URI"] = "http://example.com/api/v1.0/insert"
	configuration["VERBOSE"] = "true"
	return configuration
}

// Get required fields in configuration file
func (service *Service) GetRequired() []string{
	return []string{
		"EMAIL",
		"TO_SERVER",
		"SERVER_URI",
		"VERBOSE",
	}
}

// Fetching username tweets
func (service *Service) Fetch(configuration map[string]string, email string) (bool, error){
	module, err := service.session.SearchModule("email_to_domain")
	if err != nil {
		return false, err
	}
	targetId, err := service.session.AddTarget("email", email)
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") || targetId == ""{
			return false, err
		}
	}

	_, err = module.SetParameter("TARGET", targetId)
	if err != nil {
		return false, err
	}

	module.Start()

	target, err := service.session.GetTarget(targetId)
	if err != nil {
		return false, err
	}

	results, err := target.GetFormatedResults("email_to_domain")

	js, err := json.Marshal(&results)
	if err != nil{
		return false, err
	}

	if strings.ToLower(configuration["TO_SERVER"]) == "true" {
		if strings.ToLower(configuration["VERBOSE"]) == "true" {
			log.Println("Prepare request '"+email+"' at '"+configuration["SERVER_URI"]+"'")
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
			log.Println("Request '"+email+"' as been sent at '"+configuration["SERVER_URI"]+"'")
		}
	}
	return true, nil
}

// Service Execution with opf routine
func (service *Service) Run() (bool, error){

	configuration, _ := godotenv.Read(service.session.Config.Common.ConfigurationService + service.Name() + "/service.conf")
	service.session.Stream.Verbose = false

	if strings.Contains(configuration["EMAIL"], ",") {
		emails := strings.Split(configuration["EMAIL"], ",")
		for _, email := range emails{
			ret, err := service.Fetch(configuration, email)
			if err != nil {
				return ret, err
			}
		}
	} else {
		ret, err := service.Fetch(configuration, configuration["EMAIL"])
		if err != nil {
			return ret, err
		}
	}
	return true, nil
}
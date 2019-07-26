// AUTHOR: 	TRISTAN GRANIER
// RESUME: 	This service get last tweets for configured username
// 			every 3 hours, you can send results to web server.
package tweets_service

import (
	"encoding/json"
	"fmt"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/operative-framework/supervisor"
	"github.com/joho/godotenv"
	"strings"
	"time"
)

type Service struct{
	supervisor.Service
	session	*session.Session
}

func GetNewService(sess *session.Session) *Service{
	return &Service{
		session: sess,
	}
}

// Service as been started every 3 hours
func (service *Service) GetHibernate() time.Duration{
	return 3 * time.Hour
}

// Service name as been set here
func (service *Service) Name() string{
	return "tweets.service"
}

// Define if service need configuration file
func (service *Service) HasConfiguration() bool{
	return true
}

// Get configuration file
func (service *Service) GetConfiguration() string{
	return "./services/tweets.service/service.conf"
}

// Get required fields in configuration file
func (service *Service) GetRequired() []string{
	return []string{
		"TWITTER",
		"LAST_TWEETS",
	}
}

// Service Execution with opf routine
func (service *Service) Run() (bool, error){

	configuration, _ := godotenv.Read(service.GetConfiguration())

	service.session.Stream.Verbose = false

	module, err := service.session.SearchModule("twitter_tweets")
	if err != nil {
		return false, err
	}
	targetId, err := service.session.AddTarget("twitter", configuration["TWITTER"])
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") || targetId == ""{
			return false, err
		}
	}

	_, err = module.SetParameter("TARGET", targetId)
	if err != nil {
		return false, err
	}

	_, err = module.SetParameter("COUNT", configuration["LAST_TWEETS"])
	if err != nil {
		return false, err
	}

	module.Start()

	target, err := service.session.GetTarget(targetId)
	if err != nil {
		return false, err
	}

	results, err := target.GetFormatedResults("twitter_tweets")

	js, err := json.Marshal(&results)

	if err != nil{
		return false, err
	}

	// Print json results for possible request to external service
	fmt.Println(string(js))
	return true, nil
}
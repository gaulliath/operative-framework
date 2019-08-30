//	AUTHOR	:	TRISTAN GRANIER
//	RESUME	:	This service get last tweets for configured username
//	TIME	:	Every 3 hours
package tweets_cron

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/graniet/operative-framework/session"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Service struct {
	session.CronJob
	session *session.Session
}

func GetNewService(sess *session.Session) *Service {
	return &Service{
		session: sess,
	}
}

// Service as been started every 3 hours
func (service *Service) GetHibernate() time.Duration {
	return 3 * time.Hour
}

// Service name as been set here
func (service *Service) Name() string {
	return "tweets.cron"
}

// Define if service need configuration file
func (service *Service) HasConfiguration() bool {
	return true
}

// Get configuration
func (service *Service) GetConfiguration() map[string]string {
	configuration := make(map[string]string)
	configuration["TWITTER"] = "username1,username2"
	configuration["LAST_TWEETS"] = "50"
	configuration["TO_SERVER"] = "false"
	configuration["SERVER_URI"] = "http://example.com/api/v1.0/insert"
	configuration["VERBOSE"] = "true"
	return configuration
}

// Get required fields in configuration file
func (service *Service) GetRequired() []string {
	return []string{
		"TWITTER",
		"LAST_TWEETS",
		"TO_SERVER",
		"SERVER_URI",
		"VERBOSE",
	}
}

// Fetching username tweets
func (service *Service) Fetch(configuration map[string]string, username string) (bool, error) {
	module, err := service.session.SearchModule("twitter.tweets")
	if err != nil {
		return false, err
	}
	targetId, err := service.session.AddTarget("twitter", username)
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") || targetId == "" {
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

	results, err := target.GetFormatedResults("twitter.tweets")

	js, err := json.Marshal(&results)
	if err != nil {
		return false, err
	}

	if strings.ToLower(configuration["TO_SERVER"]) == "true" {
		if strings.ToLower(configuration["VERBOSE"]) == "true" {
			log.Println("Prepare request '" + username + "' at '" + configuration["SERVER_URI"] + "'")
		}
		req, err := http.NewRequest("POST", configuration["SERVER_URI"], bytes.NewBuffer(js))
		if err != nil {
			return false, errors.New("Can't make request to '" + configuration["SERVER_URI"] + "'")
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		if strings.ToLower(configuration["VERBOSE"]) == "true" {
			log.Println("Request '" + username + "' as been sent at '" + configuration["SERVER_URI"] + "'")
		}
	}
	return true, nil
}

// Service Execution with opf routine
func (service *Service) Run() (bool, error) {

	configuration, _ := godotenv.Read(service.session.Config.Common.ConfigurationJobs + service.Name() + "/cron.conf")
	service.session.Stream.Verbose = false

	if strings.Contains(configuration["TWITTER"], ",") {
		usernames := strings.Split(configuration["TWITTER"], ",")
		for _, username := range usernames {
			ret, err := service.Fetch(configuration, username)
			if err != nil {
				return ret, err
			}
		}
	} else {
		ret, err := service.Fetch(configuration, configuration["TWITTER"])
		if err != nil {
			return ret, err
		}
	}
	return true, nil
}

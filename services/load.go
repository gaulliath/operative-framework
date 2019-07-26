package services

import (
	"github.com/graniet/operative-framework/services/tweets.service"
	"github.com/graniet/operative-framework/supervisor"
	"time"
)

func Load(sup *supervisor.Supervisor) {
	// Loading tweets.service
	sup.Services = append(sup.Services, supervisor.Listener{
		Service:       tweets_service.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})
}
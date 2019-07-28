package services

import (
	"github.com/graniet/operative-framework/services/email_to_domain.service"
	"github.com/graniet/operative-framework/services/pastebin.service"
	"github.com/graniet/operative-framework/services/societe_com.service"
	"github.com/graniet/operative-framework/services/tweets.service"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/operative-framework/supervisor"
	"time"
)

func Load(sup *supervisor.Supervisor) {
	// Loading tweets.service
	sup.Services = append(sup.Services, session.Listener{
		Service:       tweets_service.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})

	// Loading pastebin.service
	sup.Services = append(sup.Services, session.Listener{
		Service: pastebin_service.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})

	// Loading email_to_domain.service
	sup.Services = append(sup.Services, session.Listener{
		Service: email_to_domain_service.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})

	// Loading societe_com.service
	sup.Services = append(sup.Services, session.Listener{
		Service: societe_com_service.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})
}
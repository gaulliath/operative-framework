package cron

import (
	"github.com/graniet/operative-framework/cron/email_to_domain.cron"
	"github.com/graniet/operative-framework/cron/pastebin.cron"
	"github.com/graniet/operative-framework/cron/societe_com.cron"
	"github.com/graniet/operative-framework/cron/tweets.cron"
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/operative-framework/supervisor"
	"time"
)

func Load(sup *supervisor.Supervisor) {
	// Loading tweets.service
	sup.Services = append(sup.Services, session.Listener{
		CronJob:       tweets_cron.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})

	// Loading pastebin.service
	sup.Services = append(sup.Services, session.Listener{
		CronJob:       pastebin_cron.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})

	// Loading email_to_domain.service
	sup.Services = append(sup.Services, session.Listener{
		CronJob:       email_to_domain_cron.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})

	// Loading societe_com.service
	sup.Services = append(sup.Services, session.Listener{
		CronJob:       societe_com_cron.GetNewService(sup.GetStandaloneSession()),
		NextExecution: time.Now(),
	})
}

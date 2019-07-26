package supervisor

import (
	"github.com/graniet/operative-framework/engine"
	"github.com/graniet/operative-framework/session"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Supervisor struct{
	Services 	[]Listener
	History		[]string
	Session     *session.Session
}

type Listener struct{
	ExecutedAt time.Time
	NextExecution time.Time
	Service Service
}

type Service interface {
	Name()	string
	Run()	(bool, error)
	GetHibernate() time.Duration
	HasConfiguration() bool
	GetConfiguration() string
	GetRequired() []string
}

func GetNewSupervisor(s *session.Session) *Supervisor{
	return &Supervisor{
		Session: s,
	}
}

func (sup *Supervisor) GetStandaloneSession() *session.Session{
	newSession := engine.New()
	newSession.PushPrompt()
	newSession.Config.Common.ConfigurationFile = sup.Session.Config.Common.ConfigurationFile
	return newSession
}

func (sup *Supervisor) AddHistory(s string) {
	sup.History = append(sup.History, s)
	return
}

func (sup *Supervisor) Launch(service Listener, routine chan int) Listener{

	if service.Service.HasConfiguration() {
		configuration, err := godotenv.Read(service.Service.GetConfiguration())
		if err != nil {
			log.Fatalln("'" + service.Service.GetConfiguration() + "' Config as been not found")
		}

		for _, validator := range service.Service.GetRequired() {
			if _, ok := configuration[validator]; !ok {
				log.Fatalln("'" + validator + "' field as required in configuration file")
			}
		}
	}

	service.ExecutedAt = time.Now()
	service.NextExecution = time.Now().Add(service.Service.GetHibernate())
	routine <- 1
	go func() {
		log.Println("execution of service:", service.Service.Name(), "at", service.ExecutedAt)
		log.Println("next execution at:", service.NextExecution)

		_, err := service.Service.Run()
		if err != nil {
			log.Fatalln(err.Error())
		}
		<-routine
	}()
	return service
}

func (sup *Supervisor) Read() {
	routine := make(chan int, 3)
	currentTime := time.Now()
	for {
		for key, listen := range sup.Services{
			currentTime = time.Now()
			if listen.NextExecution.Before(currentTime){
				sup.Services[key] = sup.Launch(listen, routine)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
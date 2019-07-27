package supervisor

import (
	"github.com/graniet/operative-framework/engine"
	"github.com/graniet/operative-framework/session"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Supervisor struct{
	Services 	[]session.Listener
	History		[]string
	Session     *session.Session
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
	newSession.Config.Common.ConfigurationService = sup.Session.Config.Common.ConfigurationService
	return newSession
}

func (sup *Supervisor) AddHistory(s string) {
	sup.History = append(sup.History, s)
	return
}

func (sup *Supervisor) Launch(service session.Listener, routine chan int) session.Listener{

	path := sup.Session.Config.Common.ConfigurationService + service.Service.Name() + "/service.conf"
	if service.Service.HasConfiguration() {
		configuration, err := godotenv.Read(path)
		if err != nil {
			log.Fatalln("'" + path + "' Config as been not found")
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
			log.Println(err.Error())
		}
		<-routine
	}()
	return service
}

func (sup *Supervisor) Configure() error {
	log.Println("Running service configuration...")
	if _, err := os.Stat(sup.Session.Config.Common.ConfigurationService); os.IsNotExist(err){
		_ = os.Mkdir(sup.Session.Config.Common.ConfigurationService, os.ModePerm)
	}
	for _, service := range sup.Services{
		if _, err := os.Stat(sup.Session.Config.Common.ConfigurationService + service.Service.Name()); os.IsNotExist(err){
			_ = os.Mkdir(sup.Session.Config.Common.ConfigurationService + service.Service.Name(), os.ModePerm)
		}

		if !service.Service.HasConfiguration(){
			continue
		}

		if _, err := os.Stat(sup.Session.Config.Common.ConfigurationService + service.Service.Name() + "/service.conf"); !os.IsNotExist(err){
			continue
		}

		// Generate a .env
		var file *os.File
		var errPath error

		file, errPath = os.OpenFile(sup.Session.Config.Common.ConfigurationService + service.Service.Name() + "/service.conf", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		if errPath != nil{
			return errPath
		}
		defer file.Close()

		// Writing parameters
		for key, param := range service.Service.GetConfiguration(){
			if param == "" {
				_, _ = file.WriteString(key + "=\n")
			} else{
				_, _ = file.WriteString(key + "=" + "\""+param+"\"\n")
			}
		}
		sup.Session.AddService(service)
	}
	return nil
}

func (sup *Supervisor) Read() {

	err := sup.Configure()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
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
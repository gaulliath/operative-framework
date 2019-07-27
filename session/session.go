package session

import (
	"github.com/chzyer/readline"
	"github.com/graniet/operative-framework/config"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type Session struct{
	Id int `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionName string `json:"session_name"`
	Information Information
	Connection Connection `json:"-" sql:"-"`
	Config config.Config
	Version string            `json:"version" sql:"-"`
	Targets []*Target   `json:"subjects" sql:"-"`
	Modules []Module          `json:"modules" sql:"-"`
	Filters []ModuleFilter `json:"filters" sql:"-"`
	Prompt *readline.Config `json:"-" sql:"-"`
	Stream Stream `json:"-" sql:"-"`
	TypeLists []string `json:"type_lists" sql:"-"`
	ServiceFolder	string `json:"home_folder"`
	Services	[]Listener
}

type Information struct{
	ApiStatus bool `json:"api_status"`
	ModuleLaunched int `json:"module_launched"`
	Event int `json:"event"`
}

type Listener struct{
	ExecutedAt time.Time `json:"executed_at"`
	NextExecution time.Time `json:"next_execution"`
	Service Service `json:"service"`
}

type Service interface {
	Name()	string
	Run()	(bool, error)
	GetHibernate() time.Duration
	HasConfiguration() bool
	GetConfiguration() map[string]string
	GetRequired() []string
}

func (i *Information) AddEvent(){
	i.Event = i.Event + 1
	return
}

func (i *Information) AddModule(){
	i.ModuleLaunched = i.ModuleLaunched + 1
	return
}

func (i *Information) SetApi(s bool){
	i.ApiStatus = s
	return
}

func (Session) TableName() string{
	return "sessions"
}


func (s *Session) GetId() int{
	return s.Id
}

func (s *Session) ListType() []string{
	return s.TypeLists
}

func (s *Session) PushType(t string){
	for _, tp := range s.TypeLists{
		if tp == t{
			return
		}
	}
	s.TypeLists = append(s.TypeLists, t)
}

func (s *Session) AddService(service Listener){
	s.Services = append(s.Services, service)
}

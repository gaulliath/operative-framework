package session

import (
	"time"

	"github.com/chzyer/readline"
	"github.com/graniet/operative-framework/config"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Session struct {
	Id                 int         `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionName        string      `json:"session_name"`
	Information        Information `json:"information"`
	Connection         Connection  `json:"-" sql:"-"`
	Client             OpfClient
	Events             Events            `json:"events"`
	Config             config.Config     `json:"config" sql:"-"`
	Version            string            `json:"version" sql:"-"`
	Targets            []*Target         `json:"subjects" sql:"-"`
	Modules            []Module          `json:"modules" sql:"-"`
	Monitors           Monitors          `json:"monitors"`
	Filters            []ModuleFilter    `json:"filters" sql:"-"`
	Prompt             *readline.Config  `json:"-" sql:"-"`
	Stream             Stream            `json:"-" sql:"-"`
	TypeLists          []string          `json:"type_lists" sql:"-"`
	ServiceFolder      string            `json:"home_folder"`
	Services           []Listener        `json:"services"`
	Alias              map[string]string `json:"-" sql:"-"`
	Interval           []*Interval       `json:"-"`
	Analytics          *Analytics        `json:"analytics"`
	LastAnalyticsModel string            `json:"analytics_model"`
	LastAnalyticsLinks string            `json:"last_analytics_links"`
}

type SessionExport struct {
	Id            int            `json:"-"`
	SessionName   string         `json:"session_name"`
	Information   Information    `json:"information"`
	Config        config.Config  `json:"config" sql:"-"`
	Version       string         `json:"version" sql:"-"`
	Targets       []*Target      `json:"subjects" sql:"-"`
	Modules       []Module       `json:"modules" sql:"-"`
	Filters       []ModuleFilter `json:"filters" sql:"-"`
	Stream        Stream         `json:"-" sql:"-"`
	TypeLists     []string       `json:"type_lists" sql:"-"`
	ServiceFolder string         `json:"home_folder"`
	Services      []Listener     `json:"services"`
}

type Information struct {
	ApiStatus      bool `json:"api_status"`
	ModuleLaunched int  `json:"module_launched"`
	Event          int  `json:"event"`
}

type Listener struct {
	ExecutedAt    time.Time `json:"executed_at"`
	NextExecution time.Time `json:"next_execution"`
	CronJob       CronJob   `json:"cron_job"`
}

type CronJob interface {
	Name() string
	Run() (bool, error)
	GetHibernate() time.Duration
	HasConfiguration() bool
	GetConfiguration() map[string]string
	GetRequired() []string
}

func (i *Information) AddEvent() {
	i.Event = i.Event + 1
	return
}

func (i *Information) AddModule() {
	i.ModuleLaunched = i.ModuleLaunched + 1
	return
}

func (i *Information) SetApi(s bool) {
	i.ApiStatus = s
	return
}

func (Session) TableName() string {
	return "sessions"
}

func (s *Session) GetId() int {
	return s.Id
}

func (s *Session) ListType() []string {
	return s.TypeLists
}

func (s *Session) PushType(t string) {
	for _, tp := range s.TypeLists {
		if tp == t {
			return
		}
	}
	s.TypeLists = append(s.TypeLists, t)
}

func (s *Session) AddService(service Listener) {
	s.Services = append(s.Services, service)
}

func (s *Session) ExportNow() SessionExport {
	export := SessionExport{
		Id:            s.Id,
		SessionName:   s.SessionName,
		Information:   s.Information,
		Config:        s.Config,
		Version:       s.Version,
		Targets:       s.Targets,
		Modules:       s.Modules,
		Filters:       s.Filters,
		Stream:        s.Stream,
		TypeLists:     s.TypeLists,
		ServiceFolder: s.ServiceFolder,
		Services:      s.Services,
	}
	return export
}

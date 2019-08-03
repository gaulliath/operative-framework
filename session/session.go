package session

import (
	"github.com/chzyer/readline"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/config"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"os"
	"time"
)

type Session struct {
	Id            int               `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionName   string            `json:"session_name"`
	Information   Information       `json:"information"`
	Connection    Connection        `json:"-" sql:"-"`
	Config        config.Config     `json:"config" sql:"-"`
	Version       string            `json:"version" sql:"-"`
	Targets       []*Target         `json:"subjects" sql:"-"`
	Modules       []Module          `json:"modules" sql:"-"`
	Filters       []ModuleFilter    `json:"filters" sql:"-"`
	Prompt        *readline.Config  `json:"-" sql:"-"`
	Stream        Stream            `json:"-" sql:"-"`
	TypeLists     []string          `json:"type_lists" sql:"-"`
	ServiceFolder string            `json:"home_folder"`
	Services      []Listener        `json:"services"`
	Alias         map[string]string `json:"-" sql:"-"`
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

func (s *Session) AddAlias(alias string, module string) {
	mod, err := s.SearchModule(module)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	lists := make(map[string]string)
	for a, als := range s.Alias {
		if als != mod.Name() {
			lists[a] = als
		}
	}

	lists[alias] = mod.Name()
	s.Alias = lists
	s.Stream.Success("Alias '" + alias + "' as created for module '" + mod.Name() + "'")
	return
}

func (s *Session) GetAlias(alias string) (string, error) {
	for a, m := range s.Alias {
		if a == alias {
			return m, nil
		}
	}
	return "", errors.New("Alias not found in current session.")
}

func (s *Session) ListAlias() {
	t := s.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"ALIAS",
		"MODULE",
	})
	for a, m := range s.Alias {
		t.AppendRow(table.Row{
			a,
			m,
		})
	}
	s.Stream.Render(t)
	return
}

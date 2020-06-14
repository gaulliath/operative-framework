package engine

import (
	"net/http"
	"os"
	"time"

	"github.com/graniet/operative-framework/config"
	"github.com/graniet/operative-framework/filters"
	"github.com/graniet/operative-framework/modules"
	"github.com/graniet/operative-framework/session"
	"github.com/jinzhu/gorm"
)

const VERSION = "1.20 (perception)"

// Generate New Session
func New() *session.Session {
	conf, err := config.ParseConfig()
	if err != nil {
		panic(err.Error())
	}
	db, err := gorm.Open(conf.Database.Driver, conf.Database.Name)
	if err != nil {
		panic(err.Error())
	}

	t := time.Now()
	timeText := t.Format("2006-01-02 15:04:05")

	s := session.Session{
		SessionName: "opf_" + timeText,
		Version:     VERSION,
		Information: session.Information{
			ApiStatus:      false,
			ModuleLaunched: 0,
			Event:          0,
		},
		Stream: session.Stream{
			Verbose: true,
		},
		Connection: session.Connection{
			ORM:        db,
			Migrations: make(map[string]interface{}),
		},
		Client:    session.GetOpfClient(),
		Config:    conf,
		Alias:     make(map[string]string),
		Analytics: &session.Analytics{},
	}
	s.Tracker.Server = &http.Server{
		Addr:    os.Getenv("TRACKING_HOST") + ":" + os.Getenv("TRACKING_PORT"),
		Handler: s.GetTrackerRouter(),
	}
	s.Stream.Sess = &s
	s.Connection.Migrate()
	modules.LoadModules(&s)
	filters.LoadFilters(&s)
	db.Create(&s)
	return &s
}

// Load Session With ID
func Load(id int) *session.Session {
	conf, err := config.ParseConfig()
	if err != nil {
		panic(err.Error())
	}
	db, err := gorm.Open(conf.Database.Driver, conf.Database.Name)
	if err != nil {
		panic(err.Error())
	}
	s := session.Session{
		Version: VERSION,
		Stream: session.Stream{
			Verbose: true,
		},
		Information: session.Information{
			ApiStatus:      false,
			ModuleLaunched: 0,
			Event:          0,
		},
		Connection: session.Connection{
			ORM:        db,
			Migrations: make(map[string]interface{}),
		},
		Client:    session.GetOpfClient(),
		Config:    conf,
		Analytics: &session.Analytics{},
	}
	s.Connection.ORM.Where(&session.Session{
		Id: id,
	}).First(&s)
	s.Stream.Sess = &s
	s.Connection.Migrate()
	modules.LoadModules(&s)
	filters.LoadFilters(&s)

	// Load targets
	var targets []*session.Target
	s.Connection.ORM.Where("session_id = ?", id).Find(&targets)
	s.Targets = targets

	// Load target results
	if len(s.Targets) > 0 {
		for k, target := range s.Targets {
			var linked []session.Linking
			target.Results = make(map[string][]*session.TargetResults)
			s.Connection.ORM.Where("session_id = ?", id).Where("target_base = ?", target.GetId()).Find(&linked)
			s.Targets[k].TargetLinked = linked
			s.Targets[k].Sess = &s
			if len(s.Modules) > 0 {
				for _, module := range s.Modules {
					var results []*session.TargetResults
					s.Connection.ORM.Where("session_id = ?", id).Where("module_name = ?", module.Name()).Where("target_id = ?", target.GetId()).Find(&results)
					if len(results) > 0 {
						target.Results[module.Name()] = results
					}
				}
			}
		}
	}
	return &s
}

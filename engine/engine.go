package engine

import (
	"github.com/graniet/operative-framework/modules"
	"github.com/graniet/operative-framework/session"
	"github.com/jinzhu/gorm"
	"time"
)

func New() *session.Session{
	db, err := gorm.Open("sqlite3", "./opf.db")
	if err != nil {
		panic(err.Error())
	}

	t := time.Now()
	timeText := t.Format("2006-01-02 15:04:05")

	s := session.Session{
		SessionName: "opf_" + timeText,
		Version: "1.00 (reborn)",
		Stream:session.Stream{
			Verbose: true,
		},
		Connection: session.Connection{
			ORM: db,
			Migrations: make(map[string]interface{}),
		},
	}
	s.Connection.Migrate()
	modules.LoadModules(&s)
	db.Create(&s)
	return &s
}

func Load(id int) *session.Session{
	db, err := gorm.Open("sqlite3", "./opf.db")
	if err != nil {
		panic(err.Error())
	}
	s := session.Session{
		Version: "1.00 (reborn)",
		Stream:session.Stream{
			Verbose: true,
		},
		Connection: session.Connection{
			ORM: db,
			Migrations: make(map[string]interface{}),
		},
	}
	s.Connection.ORM.Where(&session.Session{
		Id: id,
	}).First(&s)
	s.Connection.Migrate()
	modules.LoadModules(&s)

	// Load targets now
	var targets []*session.Target
	s.Connection.ORM.Where("session_id = ?", id).Find(&targets)
	s.Targets = targets


	// Load target result now
	if len(s.Targets) > 0{
		for _, target := range s.Targets{
			target.Results = make(map[string][]session.TargetResults)
			for _, module := range s.Modules{
				var results []session.TargetResults
				s.Connection.ORM.Where("session_id = ?", id).Where("module_name = ?", module.Name()).Where("target_id = ?", target.GetId()).Find(&results)
				if len(results) > 0 {
					target.Results[module.Name()] = results
				}
			}
		}
	}

	return &s
}

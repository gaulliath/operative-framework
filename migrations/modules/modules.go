package modules

import (
	"github.com/graniet/operative-framework/session"
	"time"
)

type Modules struct{
	ModuleId int `gorm:"primary_key:yes;column:module_id"`
	ModuleName string
	ModuleDescription string
	ModuleAuthor string
	ModuleType string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func PutMigration(sess *session.Session){
	sess.Connection.Migrations["modules"] = Modules{}
}

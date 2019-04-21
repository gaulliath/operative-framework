package migrations

import (
	"github.com/graniet/operative-framework/migrations/modules"
	"github.com/graniet/operative-framework/migrations/targets"
	"github.com/graniet/operative-framework/session"
)

func LoadMigrations(sess *session.Session){
	modules.PutMigration(sess)
	targets.PutMigration(sess)
}

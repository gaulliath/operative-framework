package targets

import (
	"github.com/graniet/operative-framework/session"
	"time"
)

type Targets struct{
	session.Target
	CreatedAt time.Time
	UpdatedAt time.Time
}

func PutMigration(sess *session.Session){
	sess.Connection.Migrations["targets"] = Targets{}
}

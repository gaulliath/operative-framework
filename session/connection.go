package session

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Connection struct {
	ConnectionInstance
	ORM *gorm.DB
	Migrations map[string]interface{}
}

type ConnectionInstance interface {
	GetORM() *gorm.DB
	GetDB() *sql.DB
	Migrate() bool
}

func (c *Connection) GetORM() *gorm.DB{
	return c.ORM
}

func (c *Connection) GetDB() *sql.DB{
	return c.ORM.DB()
}

func (c *Connection) LoadMigration(){
	c.Migrations["sessions"] = Session{}
	c.Migrations["targets"] = Target{}
	c.Migrations["target_links"] = Linking{}
	c.Migrations["target_results"] = TargetResults{}
}

func (c *Connection) Migrate() bool{
	c.LoadMigration()
	for _, migration := range c.Migrations{
		//if c.ORM.HasTable(migration){
		//	c.ORM.DropTable(migration)
		//}
		c.ORM.AutoMigrate(migration)
	}
	return true
}
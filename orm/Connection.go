package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Connection struct {
}

func (c Connection) Connect() *gorm.DB {

	db, err := gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		panic(err)
	}
	return db
}

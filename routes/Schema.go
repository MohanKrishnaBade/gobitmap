package routes

import (
	"github.com/gobitmap/orm"
)

func CreateSchema() {
	db := orm.Connection{}.Connect()
	tables := []orm.Model{
		&orm.User{},
	}
	for _, v := range tables {
		v.CreateTable(db)
	}
	defer db.Close()
}

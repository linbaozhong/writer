package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var db *xorm.Engine

func init() {
	var err error
	db, err = xorm.NewEngine("mysql", "root:goldapple@127.0.0.1:3306/writer?charset=utf8")

	db.ShowDebug = true
	db.ShowErr = true
	db.ShowInfo = true
	db.ShowWarn = true

	if err != nil {
		trace(err)
	}
}

func trace(arg ...interface{}) {
	//db.LogDebug(arg)
}

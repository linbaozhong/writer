package models

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	Role_SuperAdmin    = 0 //超级管理员
	Role_Administrator = 1 //管理员
	Role_Founder       = 2 //创建者
	Role_Author        = 3 //作者
	Role_Reader        = 4 //读者

	Open_Alipay = 0 //支付宝用户
	Open_QQ     = 1 //qq用户
	Open_Weibo  = 2 //微博用户

	Private  = 0 //私有的
	Internal = 1 //内部的
	Public   = 2 //公开的
	Free     = 3 //自由的

	UnDeleted = 0 //未删除
	Deleted   = 1 //删除

	UnLocked = 0 //未锁定
	Locked   = 1 //锁定
)

var db *xorm.Engine

func init() {
	var err error
	db, err = xorm.NewEngine("mysql", "root:goldapple@127.0.0.1:3306/writer?charset=utf8")

	if err != nil {
		trace(err)
	}
}

func trace(arg ...interface{}) {
	beego.Trace(arg)
}

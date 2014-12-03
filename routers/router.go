package routers

import (
	"writer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.Account{})
	//beego.AutoRouter(&controllers.Account{})
	beego.Trace("haha")
}

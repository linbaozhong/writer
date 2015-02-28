package routers

import (
	"writer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// 屏蔽路由大小写敏感
	beego.RouterCaseSensitive = false

	home := &controllers.Home{}
	beego.Router("/", home, "get:Get")
	beego.Router("/r/:parentId", home, "get:Read")
	beego.Router("/r/:parentId/:articleId", home, "get:Read")
	beego.AutoRouter(home)

	conn := &controllers.Connect{}
	beego.Router("/connect/qq_error/:msg", conn)
	beego.AutoRouter(conn)

	beego.Router("/profile", &controllers.Profile{})

	act := &controllers.Account{}
	beego.Router("/login", act, "get:Login")
	beego.Router("/signin", act, "post:SignIn")
	beego.Router("/signout", act, "post:SignOut")
	beego.Router("/signup", act, "get:SignUp")
	beego.Router("/passwordreset", act, "get:PasswordReset")

	article := &controllers.Article{}
	beego.Router("/w", article, "get:Write")
	beego.AutoRouter(article)
}

package main

import (
	"github.com/astaxie/beego"

	"writer/models"
	_ "writer/routers"
	"zouzhe/utils"
)

func main() {
	beego.Run()
}

func init() {
	//模板后缀
	beego.AddTemplateExt(".html")
	//静态目录
	beego.SetStaticPath("/htm", "html")
}

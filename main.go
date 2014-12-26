package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"writer/models"
	_ "writer/routers"
	"zouzhe/utils"
)

func main() {
	fmt.Println(models)
	beego.Run()
}

func init() {
	//模板后缀
	beego.AddTemplateExt(".html")
	//静态目录
	beego.SetStaticPath("/htm", "html")
}

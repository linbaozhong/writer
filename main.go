package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"os"
	"time"
	_ "writer/routers"
)

const (
	APP_VER = "0.1.14.1105"
)

func main() {
	beego.Info(beego.AppName, APP_VER)
	beego.Run()
}

func init() {
	//beego.SetLevel(beego.LevelInformational)
	os.Mkdir("./log", os.ModePerm)
	//日志文件名
	beego.BeeLogger.SetLogger("file", `{"filename": "log/log.log"}`)

	initStaticPath()
	initTemplateExt()
	initFuncMap()
}

// 初始化模板函数
func initFuncMap() {
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.AddFuncMap("loadtimes", loadtimes)
}

// 初始化静态目录
func initStaticPath() {
	beego.SetStaticPath("/htm", "html")

}

// 初始化模板后缀
func initTemplateExt() {
	beego.AddTemplateExt(".html")
}

// 引用自beego官网
func loadtimes(t time.Time) int {
	return int(time.Now().Sub(t).Nanoseconds() / 1e6)
}

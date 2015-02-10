package controllers

import (
	//"strings"
	//"fmt"
	//"github.com/astaxie/beego"
	"zouzhe/utils"
)

type Home struct {
	Front
}

func (this *Home) Get() {
	//this.LayoutSections["scripts"] = strings.ToLower(this.controllerName) + "/_index.html"
	this.Data["accoundId"] = utils.Str2int64(this.Ctx.GetCookie("_snow_id"))

	this.setTplNames("index")
}

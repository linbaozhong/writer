package controllers

import (
//"strings"
//"fmt"
//"github.com/astaxie/beego"
)

type Home struct {
	Front
}

func (this *Home) Get() {
	//this.LayoutSections["scripts"] = strings.ToLower(this.controllerName) + "/_index.html"
	this.setTplNames("index")
}

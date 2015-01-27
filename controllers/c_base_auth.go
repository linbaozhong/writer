package controllers

import (
	"zouzhe/utils"
)

type Auth struct {
	Base
}

func (this *Auth) Prepare() {
	this.Base.Prepare()
	// 检查当前用户是否合法用户
	if !this.allowRequest() {
		if this.IsAjax() {
			this.renderJson(utils.JsonMessage(false, "", "无效用户,请登录……"))
			this.end()
		} else {
			// 跳转到错误页
			this.Redirect("/login?returnurl="+this.Ctx.Request.URL.String(), 302)
			this.end()
		}
	}
	this.Layout = "_frontLayout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head"] = "_head.html"
	this.LayoutSections["Header"] = "_header.html"
	this.LayoutSections["Login"] = "_login.html"
	this.LayoutSections["Footer"] = "_footer.html"
	this.LayoutSections["Scripts"] = "_scripts.html"
}

func (this *Auth) Finish() {
	this.trace(this.Lang)
}

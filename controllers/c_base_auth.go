package controllers

import (
	//"strings"
	//"writer/models"
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
}

func (this *Auth) Finish() {
	this.trace(this.Lang)
}

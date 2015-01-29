package controllers

import (
	"strings"
	"writer/models"
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

// 修改locked、deleted等状态字段的公共方法
func (this *Auth) status(table string, action string) {
	if id, err := this.GetInt64("id"); err == nil {
		var e error
		switch strings.ToLower(action) {
		case "lock":
			e = models.Lock(table, id)
		case "unlock":
			e = models.UnLock(table, id)
		case "delete":
			e = models.Delete(table, id)
		case "undelete":
			e = models.UnDelete(table, id)
		}
		if e == nil {
			this.renderJson(utils.JsonMessage(true, "", ""))
		} else {
			this.renderJson(utils.JsonMessage(false, "", e.Error()))
		}
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}

}

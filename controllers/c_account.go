package controllers

import (
	"fmt"
	"zouzhe/models"
	"zouzhe/utils"

	"github.com/astaxie/beego"
)

type Account struct {
	Base
}

func (this *Account) Prepare() {
	this.Base.Prepare()
	// 检查当前用户是否合法用户
	// if !this.allowRequest() {
	// 	if this.IsAjax() {
	// 		this.renderJson(utils.JsonMessage(false, "", "无效用户,请登录……"))
	// 		this.end()
	// 	} else {
	// 		// 跳转到错误页
	// 		this.Redirect("/login?returnurl="+this.Ctx.Request.URL.String(), 302)
	// 		this.end()
	// 	}
	// }
	this.Layout = "_noneLayout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head"] = "_head.html"
	this.LayoutSections["Header"] = "_noneHeader.html"
	this.LayoutSections["Login"] = ""
	this.LayoutSections["Footer"] = "_footer.html"
	this.LayoutSections["Scripts"] = "_scripts.html"
}

/*
* 新账户
 */
func (this *Account) New() {

	oa := &models.Accounts{OpenId: this.GetString("openId"), OpenFrom: this.GetString("openFrom")}

	//是否已经存在
	if has, err := oa.Exists(); has {
		//fmt.Println(oa)
	} else {

		oa.AccessToken = this.GetString("accessToken")
		oa.Avatar_1 = this.GetString("avatar")
		oa.NickName = this.GetString("nickName")

		this.extend(oa)

		if _, err, _ = oa.Post(); err == nil {
			this.trace(oa)
		} else {
			this.trace(err)
		}

	}
	//登录状态存入cookie，缺省时间是1年：365*24*60*60
	var cookieDuration interface{}
	always := false
	if always {
		cookieDuration, _ = beego.AppConfig.Int("CookieDuration")

	} else {
		cookieDuration = ""
	}

	this.Ctx.SetCookie(beego.AppConfig.String("CookieName"),
		utils.CookieEncode(fmt.Sprintf("%d|%s|%s", oa.Id, oa.NickName, oa.Avatar_1)),
		cookieDuration, "/")

	this.renderJson(utils.JsonData(true, "", oa))
}

// 注册
func (this *Account) SignUp() {
	_m_account := new(models.Accounts)
	_m_account.LoginName = this.GetString("loginName")
	_m_account.Password = this.GetString("password")

	this.Data["sign"] = _m_account
	this.Data["auto"] = this.getCheckboxBool("auto")
	this.setTplNames()
}

// 签入
func (this *Account) SignIn() {
	_m_account := new(models.Accounts)
	_m_account.LoginName = this.GetString("loginName")
	_m_account.Password = this.GetString("password")
	_m_account.Gender = this.getCheckboxInt("auto")

	this.renderJson(utils.JsonData(true, "", _m_account))
}

// 签出
func (this *Account) SignOut() {
	this.loginOut()
	this.renderJson(utils.JsonMessage(true, "", ""))
}

// 密码重置
func (this *Account) PasswordReset() {
	this.setTplNames()
}

// 独立登录页面
func (this *Account) Login() {
	this.Data["returnUrl"] = this.GetString("returnurl")
	this.setTplNames()
}

/*
* 实现第三方账号登录的相关逻辑
 */
package controllers

import (
	"errors"
	"strconv"
	"strings"
	"zouzhe/models"
	"zouzhe/utils"

	"github.com/astaxie/beego/httplib"
)

type Connect struct {
	Base
}

type OpenSign struct {
	Ret      string //返回码 qq=0表示成功
	Msg      string //错误时的消息提示
	Id       int64  //本站账户id
	From     string //第三方标识
	OpenId   string //账户id
	Token    string //access_token
	Refresh  string //refresh_token
	NickName string //昵称
	Gender   string //性别
	Avatar_1 string //40*40的qq头像
	Avatar_2 string //100*100的qq头像
}

var (
	appid    = appconf("qq::appid")
	callback = appconf("qq::callback")
	appkey   = appconf("qq::appkey")
	//---只读取用户信息
	scope = ""
)

/*
* QQ登录
 */
func (this *Connect) QQ_Login() {
	//---生成唯一随机串防止csrf攻击
	state := this.XsrfToken()
	//---将随机串存入Session
	this.SetSession("state", state)
	//---登录的url
	url := appconf("qq::auth") + "?response_type=code&client_id=" + appid + "&redirect_uri=" + callback + "&state=" + state + "&scope=" + scope

	//fmt.Println(url)
	//---
	this.Redirect(url, 302)

	//fmt.Println("after redirect ...")
}

/*
* QQ登录回调
 */
func (this *Connect) QQ_Callback() {

	//---检查是否有接口调用错误
	if this.GetString("msg") != "" {
		//---跳转至错误页
		this.Redirect(this.UrlFor("Connect.Connect_Error", ":msg", this.GetString("msg")), 302)
		return
	}

	//---验证state防止csrf攻击
	if this.GetString("state") != this.GetSession("state").(string) {
		//---跳转至错误页
		this.Redirect(this.UrlFor("Connect.Connect_Error", ":msg", ""), 302)
		return
	}

	//---opensign第三方账户信息
	_account := new(OpenSign)
	_account.From = "qq"

	//---创建读取token的请求
	req := httplib.Get(appconf("qq::token"))
	req.Param("grant_type", "authorization_code")
	req.Param("client_id", appid)
	req.Param("client_secret", appkey)
	req.Param("redirect_uri", callback)
	req.Param("code", this.GetString("code"))
	req.Param("state", this.GetString("state"))

	//---读取返回的内容
	rep, err := req.String()

	if err == nil {

		jmap := make(map[string]string)
		//---解析返回的内容
		param := strings.Split(rep, "&")

		for _, item := range param {
			pos := strings.Index(item, "=")
			jmap[item[:pos]] = item[pos+1:]
		}
		//--- access_token
		_account.Token = jmap["access_token"]
		//--- refresh_token
		_account.Refresh = jmap["refresh_token"]

	} else {
		//---跳转至错误页
		this.Redirect(this.UrlFor("Connect.Connect_Error", ":msg", err.Error()), 302)
		return
	}

	//---创建读取openid的请求
	if qq_openid(_account) != nil {
		//---跳转至错误页
		this.Redirect(this.UrlFor("Connect.Connect_Error", ":msg", err.Error()), 302)
		return
	}

	//---创建读取userinfo的请求
	if qq_userinfo(_account) != nil {
		//---跳转至错误页
		this.Redirect(this.UrlFor("Connect.Connect_Error", ":msg", err.Error()), 302)
		return
	}

	// 写入cookie
	this.cookie("openid", _account.OpenId)
	this.cookie("token", _account.Token)
	this.cookie("nickname", _account.NickName)
	this.cookie("avatar", _account.Avatar_1)

	this.trace(_account)
	//
	this.Data["sign"] = _account

	this.setTplNames("callback")

}

/*
* QQ登录获取openid
 */
func qq_openid(act *OpenSign) (err error) {
	//---创建读取openid的请求
	req := httplib.Get(appconf("qq::openid"))
	req.Param("access_token", act.Token)

	//---读取返回的内容
	rep, err := req.String()

	if err == nil {

		jmap := make(map[string]interface{})
		//---解析返回的内容,检查如果包含callback,读取openid
		if strings.Contains(rep, "callback") {
			jmap = utils.JsonString2map(rep[strings.Index(rep, "(")+1 : strings.LastIndex(rep, ")")])
		}

		if len(jmap) > 0 {
			//--- openid
			act.OpenId = utils.Interface2str(jmap["openid"])
		} else {
			err = errors.New("return value is empty")
		}
	}
	return
}

/*
* QQ登录用户信息
 */
func qq_userinfo(act *OpenSign) (err error) {
	//---创建读取userinfo的请求
	req := httplib.Get(appconf("qq::userinfo"))
	req.Param("access_token", act.Token)
	req.Param("openid", act.OpenId)
	req.Param("oauth_consumer_key", appid)
	req.Debug(true)
	//---读取返回的内容
	rep, err := req.String()

	if err == nil {
		//---解析返回的内容,检查如果包含callback,读取openid
		jmap := utils.JsonString2map(rep)

		if len(jmap) > 0 {

			act.Ret = utils.Interface2str(jmap["ret"])
			//--- 检查返回码是否正确
			if act.Ret == "0" {
				act.NickName = utils.Interface2str(jmap["nickname"])
				act.Gender = utils.Interface2str(jmap["gender"])
				act.Avatar_1 = utils.Interface2str(jmap["figureurl_qq_1"])
				act.Avatar_2 = utils.Interface2str(jmap["figureurl_qq_2"])
			} else {
				act.Msg = utils.Interface2str(jmap["msg"])
				err = errors.New(act.Msg)
			}
		} else {
			err = errors.New("return value is empty")
		}

	}
	return
}

/*
* QQ登录access_token续期
 */
func qq_refresh(act *OpenSign) (err error) {
	//---创建读取token的请求
	req := httplib.Get(appconf("qq::token"))
	req.Param("grant_type", "refresh_token")
	req.Param("client_secret", appkey)
	req.Param("client_id", appid)
	req.Param("refresh_token", act.Refresh)

	//---读取返回的内容
	rep, err := req.String()

	if err == nil && strings.Contains(rep, "access_token") && strings.Contains(rep, "refresh_token") {

		jmap := make(map[string]string)
		//---解析返回的内容
		param := strings.Split(rep, "&")

		for _, item := range param {
			pos := strings.Index(item, "=")
			jmap[item[:pos]] = item[pos+1:]
		}
		//--- access_token
		act.Token = jmap["access_token"]
		//--- refresh_token
		act.Refresh = jmap["refresh_token"]
	}
	return
}

/*
* 登录错误地址
 */
func (this *Connect) Connect_Error() {
	this.trace(this.GetString("msg"))
	this.Data["message"] = this.GetString("msg")
}

/*
* 记录登录历史
* 1.检查该账户是否存在，如果不存在，创建之，反之，记录登录日期和ip地址
* 2.检查是否需要进行授权续期，如需要，续期
 */
func (this *Connect) SignTrace() {

	_account := new(OpenSign)
	_account.Id, _ = strconv.ParseInt(this.Ctx.GetCookie("_snow_id"), 10, 64)
	_account.From = this.GetString("from")
	_account.Gender = this.GetString("gender")
	_account.NickName = this.GetString("nickName")
	_account.OpenId = this.GetString("openId")
	_account.Token = this.GetString("token")
	_account.Refresh = this.GetString("refresh")
	_account.Avatar_1 = this.GetString("avatar_1")
	_account.Avatar_2 = this.GetString("avatar_2")

	// 账号id和第三方账号id均为空，视为无效的请求
	if _account.Id == 0 && _account.OpenId == "" {
		this.trace("无效的账户信息")
		this.renderJson(utils.JsonData(false, "", errors.New("无效的账户信息")))
		return
	}

	var _m_account *models.Accounts
	// 如果_account.Id>0 检查该账户是否存在
	if _account.Id > 0 {
		_m_account = &models.Accounts{Id: _account.Id}
	} else {
		_m_account = &models.Accounts{OpenId: _account.OpenId, OpenFrom: _account.From}
	}
	// 账户是否存在
	has, err := _m_account.Exists()
	if err != nil {
		this.trace(err.Error())
		this.renderJson(utils.JsonData(false, "", err))
		return
	}
	// 如果账户存在
	if has {
		this.trace("记录登录日志")

		// 记录登录日志
		_m_log := new(models.LoginLog)
		_m_log.AccountId = _m_account.Id

		this.extend(_m_log)
		// 写登录日志
		_, err = _m_log.Post()

		// 检查access_token是否需要续期(有效期一般是3个月)
		if _m_account.RefreshToken != "" && _m_account.AccessToken != "" {
			_needRefresh := (utils.Msec2Time(_m_log.Updated).Sub(utils.Msec2Time(_m_account.Updated)).Hours() > 24*40*2)
			// 2.续期
			if _needRefresh {
				this.trace("access_token续期")
				err = qq_refresh(_account)
				if err != nil {
					// 写入cookie
					this.cookie("token", _account.Token)
					// 更新access_token and refresh_token and Updated
					_m_account.AccessToken = _account.Token
					_m_account.RefreshToken = _account.Refresh
					this.extend(_m_account)
					_m_account.RefreshAccessToken()
				}
			}
		}
	} else {
		this.trace("创建新的账户")
		// 反之，创建新账户
		_m_account.OpenFrom = _account.From
		if _account.Gender == "男" {
			_m_account.Gender = 1
		}
		_m_account.NickName = _account.NickName
		_m_account.AccessToken = _account.Token
		_m_account.RefreshToken = _account.Refresh
		_m_account.Avatar_1 = _account.Avatar_1
		_m_account.Avatar_2 = _account.Avatar_2

		this.extend(_m_account)
		// errs:记录返回的数据校验错误
		_, err, errs := _m_account.Post()

		if err != nil {
			this.renderJson(utils.JsonData(false, "", errs))
			return
		}
	}

	this.currentUser.Id = _m_account.Id
	this.currentUser.From = _m_account.OpenFrom

	// 保存登录状态
	this.loginIn(_m_account.Id, _m_account.OpenFrom)

	this.renderJson(utils.JsonMessage(true, "", ""))
}

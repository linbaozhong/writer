package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"
	"html/template"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zouzhe/models"
	"zouzhe/utils"
)

type langType struct {
	Lang, Name string
}
type Base struct {
	beego.Controller
	i18n.Locale
	page        models.Page
	currentUser *models.Current
}

//全部单词字符包括中文
const sub = "\\w\u4e00-\u9fa5"

var (
	siteDomain string //网站域
	langTypes  []*langType
)

func init() {
	siteDomain = appconf("site::domain")

	// 引用beego官网代码
	langs := strings.Split(appconf("lang::types"), "|")
	names := strings.Split(appconf("lang::names"), "|")
	langTypes = make([]*langType, 0, len(langs))
	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Lang: v,
			Name: names[i],
		})
	}

	for _, lang := range langs {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}

//读取配置
func appconf(key string) string {
	return beego.AppConfig.String(key)
}

func (this *Base) Prepare() {
	this.currentUser = new(models.Current)

	this.initPage()
}

//
func (this *Base) initPage() {

	this.Data["PageStartTime"] = time.Now()

	// Redirect to make URL clean.
	if this.setLangVer() {
		i := strings.Index(this.Ctx.Request.RequestURI, "?")
		this.Redirect(this.Ctx.Request.RequestURI[:i], 302)
		return
	}
	//this.page.Author = this.lang("appauthor")
	//this.page.Company = this.lang("appcompany")
	//this.page.Copyright = this.lang("appcopyright")
	//this.page.Description = this.lang("appdescription")
	//this.page.Domain = this.lang("appdomain")
	//this.page.Keywords = this.lang("appkeywords")
	//this.page.SiteName = this.lang("appsitename")
	//this.page.Title = this.lang("apptitle")
	//this.page.Product = this.lang("appproduct")
	//this.page.Version = this.lang("appversion")

	//this.Data["page"] = &this.page
}

/*
* setLangVer 设置网址语言版本.引用beego官网
 */
func (this *Base) setLangVer() bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	lang := this.Input().Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	if !hasCookie {
		this.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	// Set language properties.
	this.Lang = lang
	this.Data["Lang"] = curLang.Lang
	this.Data["CurLang"] = curLang.Name
	this.Data["RestLangs"] = restLangs

	return isNeedRedir
}

//公共字段
func (this *Base) extend(bean interface{}) {
	d := reflect.Indirect(reflect.ValueOf(bean))
	if v := d.FieldByName("Updator"); v.IsValid() && v.Int() == 0 && this.currentUser != nil {
		v.SetInt(this.currentUser.Id)
	}
	if v := d.FieldByName("Updated"); v.IsValid() && v.Int() == 0 {
		v.SetInt(utils.Millisecond(time.Now()))
	}
	if v := d.FieldByName("Ip"); v.IsValid() && v.String() == "" {
		v.SetString(this.Ctx.Input.IP())
	}
}

//公共字段
func (this *Base) extendEx(bean interface{}) {
	d := reflect.Indirect(reflect.ValueOf(bean))

	if v := d.FieldByName("Updator"); v.IsValid() && v.Int() == 0 && this.currentUser != nil {
		v.SetInt(this.currentUser.Id)
	}
	if v := d.FieldByName("Updated"); v.IsValid() && v.Int() == 0 {
		v.SetInt(utils.Millisecond(time.Now()))
	}
	if v := d.FieldByName("Creator"); v.IsValid() && v.Int() == 0 && this.currentUser != nil {
		v.SetInt(this.currentUser.Id)
	}
	if v := d.FieldByName("Created"); v.IsValid() && v.Int() == 0 {
		v.SetInt(utils.Millisecond(time.Now()))
	}
	if v := d.FieldByName("Ip"); v.IsValid() && v.String() == "" {
		v.SetString(this.Ctx.Input.IP())
	}
}

/*
* 读取request数据填充struct
 */
func (this *Base) fillModel(bean interface{}) {
	// 获取接口指向的实例
	v := reflect.ValueOf(bean).Elem()
	// 实例的类型
	t := v.Type()
	// 判断实例的类型是否struct
	if t.Kind() == reflect.Struct {
		//遍历struct的field
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if d := v.FieldByName(f.Name); d.IsValid() {
				// 根据field数据类型,读取request数据并赋值
				switch f.Type.Kind() {
				//整形
				case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
					n, _ := this.GetInt64(strings.ToLower(f.Name))
					d.SetInt(int64(n))
				//字符串
				case reflect.String:
					d.SetString(this.GetString(strings.ToLower(f.Name)))
					//其他类型，以此实现
					//case reflect.Bool:
				}
			}
		}
	}
}

/*
* 输出 Json 格式数据
 */
func (this *Base) outputJson(data interface{}, err error) {
	if err == nil {
		this.renderJson(utils.JsonData(true, "", data))
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}
}

//
func (this *Base) setJsonData(data interface{}) {
	//操作成功，清除token
	if resp := reflect.Indirect(reflect.ValueOf(data)); resp.FieldByName("Ok").Bool() {
		this.XsrfToken()
	}
	this.Data["json"] = data
}

//返回json响应格式
func (this *Base) renderJson(data interface{}) {
	this.setJsonData(data)
	this.ServeJson()
}

//返回jsonp响应
func (this *Base) renderJsonp(data interface{}) {
	this.setJsonData(data)
	this.ServeJsonp()
}

//返回html字符串格式响应
func (this *Base) serveString(arg string) {
	this.Ctx.Output.Body([]byte(arg))
}

//响应签名丢失错误
func (this *Base) renderLoseToken() {
	data := utils.JsonMessage(false, "invalidFormToken", "invalidFormToken")
	this.renderJson(data)
}

//是否外链
func (this *Base) isOutLink() bool {
	host, err := url.Parse(this.Ctx.Request.Referer())
	if err != nil {
		return true
	}
	return this.Ctx.Request.Host != host.Host
}

// 渲染字符串模板
func (this *Base) renderTemplateString(tplString string) error {
	// 读取当前控制器和方法名称
	a, c := this.GetControllerAndAction()

	// 查找模板是否已经存在
	name := "/" + a + "/" + c
	t := beego.BeeTemplates[name]

	if t == nil {
		// 解析传入的字符串
		var err error
		t, err = template.New(name).Parse(tplString)
		if err != nil {
			return err
		}

		t = t.Delims(beego.TemplateLeft, beego.TemplateRight)

		beego.BeeTemplates[name] = t
	}

	// 渲染并输出
	return t.Execute(this.Ctx.ResponseWriter, this.Data)
}

////文件服务
//func (this *Base) serverFile(file, filename string) {
//	file = filepath.Join(".", file)

//	//友好文件名
//	if len(filename) == 0 {
//		filename = utils.UrlEncode(filepath.Base(file))
//	}

//	this.Ctx.ResponseWriter.Header().Set("Content-Description", "File Transfer")
//	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream;charset=UTF-8")
//	this.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+filename)
//	this.Ctx.ResponseWriter.Header().Set("Content-Transfer-Encoding", "binary")
//	this.Ctx.ResponseWriter.Header().Set("Expires", "0")
//	this.Ctx.ResponseWriter.Header().Set("Cache-Control", "must-revalidate")
//	this.Ctx.ResponseWriter.Header().Set("Pragma", "public")

//	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, file)
//	this.end()
//}

//获取URL参数
func (this *Base) getParamsInt64(key string) (int64, error) {
	i64, err := strconv.ParseInt(this.getParamsString(key), 10, 64)
	return i64, err
}

func (this *Base) getParamsInt(key string) (int, error) {
	i64, err := this.getParamsInt64(key)
	return int(i64), err
}

func (this *Base) getParamsString(key string) string {
	return this.Ctx.Input.Param(key)
}

//验证合法用户
// func (this *Base) validUser() (*models.Current, bool) {

// 	this.currentUser.Id , _ = utils.Str2int64(this.Ctx.GetCookie("_snow_id"))

// 	if this.currentUser.Id == 0 {
// 		return nil, false
// 	}
// 	// 来自第三方平台的账户
// 	this.currentUser.From = this.Ctx.GetCookie("from")
// 	//
// 	if this._sonw_token(this.currentUser.Id,this.currentUser.From )==this.Ctx.GetCookie("_snow_key") {

// 	}

// 	return this.currentUser, true
// }

//允许新的请求，数据通用字段初始信息，附带验证用户是否合法(err)，
func (this *Base) allowRequest() bool {

	this.trace(this.Ctx.GetCookie("_snow_id"), this.Ctx.GetCookie("from"))
	this.currentUser.Id, _ = strconv.ParseInt(this.Ctx.GetCookie("_snow_id"), 10, 64)

	if this.currentUser.Id == 0 {
		return false
	}
	// 来自第三方平台的账户
	this.currentUser.From = this.Ctx.GetCookie("from")
	//
	return this._sonw_token(this.currentUser.Id, this.currentUser.From) == this.Ctx.GetCookie("_snow_token")
}

// //读取登录用户的Cookie信息
// func (this *Base) GetCurrentUser(cookie string) (currentuser *models.Current) {
// 	currentuser = new(models.Current)

// 	cookie = utils.CookieDecode(cookie)

// 	//拆分cookie
// 	curr := strings.Split(cookie, "|")
// 	if len(curr) > 0 {
// 		currentuser.Id, _ = utils.Str2int64(curr[0]) //strconv.ParseInt(curr[0], 10, 0)
// 	}
// 	if len(curr) > 1 {
// 		currentuser.Name = curr[1]
// 	}
// 	if len(curr) > 2 {
// 		currentuser.Avatar = curr[2]
// 	}
// 	if len(curr) > 3 {
// 		currentuser.Role = curr[3]
// 	}
// 	return
// }

/*
* 返回form表单中checkbox的状态值的bool形式
 */
func (this *Base) getCheckboxBool(key string) bool {
	return strings.ToLower(this.GetString(key)) == "on"
}

/*
* 返回form表单中checkbox的状态值的int形式
 */
func (this *Base) getCheckboxInt(key string) int {
	if this.getCheckboxBool(key) {
		return 1
	} else {
		return 0
	}
}

/*
* 解析首尾#中的字符串
 */
func (this *Base) parseSharp(str string) []string {
	//匹配字符串
	p := fmt.Sprintf("#([%s]+)#[^%s]*", sub, sub)
	return this.parseString(str, p)
}

/*
* 解析首字符@尾字符是空格中的子串
 */
func (this *Base) parseAite(str string) []string {
	//匹配字符串
	p := fmt.Sprintf("@([%s]+)", sub)
	return this.parseString(str, p)
}

/*
* 解析指定首尾字符中的字符串
 */
func (this *Base) parseString(str, p string) []string {
	//正则
	re := regexp.MustCompile(p)

	result := make([]string, 0)
	//查找子串
	for _, tags := range re.FindAllStringSubmatch(str, -1) {
		for _, tag := range tags[1:] {
			if tag != "" {
				result = append(result, tag)
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

/*
* xsrf过滤
 */
func (this *Base) checkXsrf() (bool, string) {
	if this.CheckXsrfCookie() {
		return true, this.XsrfToken()
	}
	return false, ""
}

////获取当前语言
//func (this *Base) lang(k string) string {
//	return utils.Lang(k, this.Ctx.Request.Header.Get("Accept-Language"))
//}

//终止服务
func (this *Base) end() {
	this.Layout = ""
	this.TplNames = ""

	this.StopRun()
}

// 写入cookie
func (this *Base) cookie(name, value string) {
	this.Ctx.SetCookie(name, value, 1<<31-1, "/", siteDomain)
}

// 写入cookie,禁止客户端读取
func (this *Base) cookieHttpOnly(name, value string) {
	this.Ctx.SetCookie(name, value, 1<<31-1, "/", siteDomain, nil, true)
}

// 设置模板文件
func (this *Base) setTplNames(name ...string) {
	c, a := this.Controller.GetControllerAndAction()

	if len(name) > 0 && name[0] != "" {
		a = name[0]
	}
	this.TplNames = strings.ToLower(fmt.Sprintf("%s/%s.html", c, a))
}

//签名
func (this *Base) _sonw_token(id int64, from string) string {
	return utils.MD5Ex(fmt.Sprintf("%d_%s", id, from))
}

// 签入
func (this *Base) loginIn(id int64, from string) {
	this.cookieHttpOnly("from", from)
	this.cookieHttpOnly("_snow_id", strconv.FormatInt(id, 10))
	this.cookie("_snow_token", this._sonw_token(id, from))
}

// 签出
func (this *Base) loginOut() {
	this.cookie("_snow_token", "")
}

/*
* 跟踪
 */
func (this *Base) trace(v ...interface{}) {
	c, a := this.Controller.GetControllerAndAction()
	beego.Trace(fmt.Sprintf("%s/%s ", c, a) + fmt.Sprintf("Info:%s", utils.Interface2str(v...)))
}

/*
* 读取数据校验错误
 */
func getValidErrors(valid *validation.Validation) []models.Error {
	errs := make([]models.Error, 0)
	for _, err := range valid.Errors {
		errs = append(errs, models.Error{Key: err.Key, Message: err.Message})
	}
	return errs
}

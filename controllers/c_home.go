package controllers

import (
	"strings"
	"writer/models"
	"zouzhe/utils"
)

type Home struct {
	Front
}

func (this *Home) Get() {
	//this.page.Title = "haa"
	this.Data["account"] = this.currentUser

	this.setTplNames("index")
}

// 阅读
func (this *Home) Read() {
	parentId, _ := this.getParamsInt64(":parentid")
	articleId, _ := this.getParamsInt64(":articleid")

	if parentId <= 0 {
		// 跳转到首页
		this.Redirect("/", 302)
		this.end()
	} else {
		this.Data["articleId"] = articleId
		this.Data["parentId"] = parentId
		this.setTplNames("read")
	}
}

// 分页拉取书籍列表
func (this *Home) Books() {
	// 是否读取当前用户的文档
	my := this.getParamsString("0")

	// 读取分页规则
	p := new(models.Pagination)

	if size, err := this.GetInt("size"); err != nil || size == 0 {
		p.Size = 20
	} else {
		p.Size = size
	}
	p.Index, _ = this.GetInt("index")
	p.Count, _ = this.GetInt("count")
	// 读取查询条件
	moreId, _ := this.GetInt64("moreId")
	tags := "'" + strings.Replace(strings.TrimSpace(this.GetString("tags")), " ", "','", -1) + "'"

	// 拉取
	a := new(models.Article)
	var as []models.Article
	var err error
	var cond string

	// 构造查询字符串
	if moreId == 0 {
		cond = "articlemore.parentId = ?"
	} else {
		cond = "articlemore.parentId = (select articleid from articlemore where id = ?)"
	}

	if my == "" {
		if len(tags) == 0 {
			as, err = a.List(p, cond, moreId)
		} else {
			as, err = a.List(p, cond+" and articles.id in (select articleId from tagarticle where tagId in (select id from tags where name in ("+tags+")))", moreId)
		}
	} else {
		if len(tags) == 0 {
			as, err = a.List(p, cond+" and articles.creator = ?", moreId, this.currentUser.Id)
		} else {
			as, err = a.List(p, cond+" and articles.id in (select articleId from tagarticle where tagId in (select id from tags where name in ("+tags+"))) and articles.creator = ?", moreId, this.currentUser.Id)
		}
	}
	this.trace(err)
	if err == nil {
		rj := new(models.ReturnJson)
		rj.Data = as
		rj.Page = p

		this.renderJson(utils.JsonData(true, "", rj))
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}

}

// 读取最常用标签
func (this *Home) Tags() {
	t := new(models.Tags)

	if ts, err := t.List(); err == nil {
		this.renderJson(utils.JsonData(true, "", ts))
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}
	//// 检查当前用户是否本地用户
	//if this.allowRequest() {

	//}
}

// 读取目录
func (this *Home) Catalog() {
	id, _ := this.GetInt64("parentId")

	if id <= 0 {
		this.renderJson(utils.JsonMessage(false, "", "参数错误"))
	}

	a := new(models.Article)
	as, err := a.Catalog(id)

	if err == nil {
		this.renderJson(utils.JsonData(true, "", as))
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}
}

// 分页读取内容
func (this *Home) Content() {
	// 读取分页规则
	p := new(models.Pagination)

	if size, err := this.GetInt("size"); err != nil || size == 0 {
		p.Size = 10000
	}
	p.Index, _ = this.GetInt("index")
	// 读取查询条件
	parentId, _ := this.GetInt64("parentId")
	articleId, _ := this.GetInt64("articleId")

	// 构造查询字符串
	cond := "articles.id >= ?"

	// 拉取
	a := new(models.Article)
	as, err := a.GetContent(p, parentId, cond, articleId)

	if err == nil {
		this.renderJson(utils.JsonData(true, "", as))
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}
}

// 读取一个内容
func (this *Home) Single() {
	// 读取查询条件
	articleId, _ := this.GetInt64("articleId")

	// 构造查询字符串
	cond := "articles.id >= ?"

	// 拉取
	a := new(models.Article)
	err := a.GetSingle(cond, articleId)

	if err == nil {
		this.renderJson(utils.JsonData(true, "", a))
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}
}

//// 搜索标签
//func (this *Home)Search(){
//	// 读取请求标签
//	tags = this.GetString('k')

//}

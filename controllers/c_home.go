package controllers

import (
	"writer/models"
	"zouzhe/utils"
)

type Home struct {
	Front
}

func (this *Home) Get() {
	this.setTplNames("index")
}

// 阅读
func (this *Home) Read() {
	this.setTplNames("read")
}

// 分页拉取书籍列表
func (this *Home) Books() {
	// 读取分页规则
	p := new(models.Pagination)

	if size, err := this.GetInt("size"); err != nil || size == 0 {
		p.Size = 20
	}
	p.Index, _ = this.GetInt("index")
	// 读取查询条件
	parentId := 0
	tags := this.GetString("tags")

	// 拉取
	a := new(models.Article)
	var as []models.Article
	var err error

	// 构造查询字符串
	cond := "parentId=?"

	if tags == "" {
		as, err = a.List(p, cond, parentId)
	} else {
		cond += " and tags=?"
		as, err = a.List(p, cond, parentId, tags)
	}

	if err == nil {
		this.renderJson(utils.JsonData(true, "", as))
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}

}

// 读取目录
func (this *Home) Catalog() {
	id, _ := this.getParamsInt64("0")

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

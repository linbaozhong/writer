package controllers

import (
	"fmt"
	"writer/models"
	"zouzhe/utils"
)

type Article struct {
	Auth
}

func (this *Article) Update() {
	article := new(models.Articles)

	if id, err := this.GetInt64("id"); err == nil {
		article.Id = id
	}
	if parentid, err := this.GetInt64("parentId"); err == nil {
		article.ParentId = parentid
	}

	article.Title = this.GetString("title")
	article.Content = this.GetString("content")
	article.Tags = this.GetString("tags")

	if position, err := this.GetInt("position"); err == nil {
		article.Position = position
	}
	//
	if article.Id > 0 {
		this.extend(article)
	} else {
		this.extendEx(article)
	}

	// 提交更新，返回结果
	if err, errs := article.Update(); err == nil {
		this.renderJson(utils.JsonData(true, "", article))
	} else {
		this.renderJson(utils.JsonData(false, "", errs))
	}

}

/*
方法：读取一条数据
参数：id
*/
func (this *Article) Get() {
	id, err := this.GetInt64("id")

	if err != nil {
		this.renderJson(utils.JsonMessage(false, "", "参数错误"))
		return
	}

	//
	a := new(models.Articles)
	a.Id = id

	if has, err := a.Get(); err == nil {
		if has {
			this.renderJson(utils.JsonData(true, "", a))
		} else {
			this.renderJson(utils.JsonMessage(false, "", "数据不存在"))
		}
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}
}

// @Title List
// @Description 分页拉取问题列表
// @Param   size  form  int  false        "每页的记录条数"
// @Success 200 {object} utils.Response
// @Failure 200 {object} utils.Response
// @router /List [post]
func (this *Article) List() {
	// 读取分页规则
	p := new(models.Pagination)

	if size, err := this.GetInt("size"); err != nil || size == 0 {
		p.Size = 20
	}
	p.Index, _ = this.GetInt("index")
	// 读取查询条件
	when := this.GetString("when")
	where := this.GetString("where")
	// 构造查询字符串
	cond := "1=1"
	if when != "" {
		cond += fmt.Sprintf(" and when='%s'", when)
	}
	if where != "" {
		cond += fmt.Sprintf(" and where='%s'", where)
	}

	// 拉取
	a := new(models.Articles)

	if as, err := a.List(cond, p); err != nil {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	} else {
		this.renderJson(utils.JsonData(true, "", as))
	}
}

// 锁定
func (this *Article) Lock() {
	this.status("Articles", "lock")
}

// 解锁
func (this *Article) UnLock() {
	this.status("Articles", "unlock")
}

// 删除
func (this *Article) Delete() {
	this.status("Articles", "delete")
}

// 恢复
func (this *Article) UnDelete() {
	this.status("Articles", "undelete")
}

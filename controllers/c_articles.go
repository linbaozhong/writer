package controllers

import (
	//"fmt"
	"writer/models"
	"zouzhe/utils"
)

type Article struct {
	Auth
}

func (this *Article) Update() {
	article := new(models.Article)

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
	a := new(models.Article)
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
	parentId, _ := this.GetInt64("parentid")
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

// @Title Position
// @Description 节点位置发生变化
// @Param parentId  form  int  false        "父级id"
// @Param position  form  int  false        "位置索引"
func (this *Article) Position() {
	a := new(models.Article)
	a.Id, _ = this.GetInt64("id")
	a.ParentId, _ = this.GetInt64("parentid")
	a.Position, _ = this.GetInt("position")

	this.extend(a)

	if a.Id <= 0 {
		this.renderJson(utils.JsonMessage(false, "id", "参数错误: id 必须 >0"))
	}

	if ok, err := a.SetPosition(); err == nil {
		if ok {
			this.renderJson(utils.JsonMessage(true, "", ""))
		} else {
			this.renderJson(utils.JsonMessage(false, "", ""))
		}
	} else {
		this.renderJson(utils.JsonMessage(false, "", err.Error()))
	}
}

// 锁定
func (this *Article) Lock() {
	this.status("lock")
}

// 解锁
func (this *Article) UnLock() {
	this.status("unlock")
}

// 删除
func (this *Article) Delete() {
	this.status("delete")
}

// 恢复
func (this *Article) UnDelete() {
	this.status("undelete")
}

func (this *Article) status(action string) {
	if id, err := this.GetInt64("id"); err == nil && id > 0 {
		a := new(models.Article)
		a.Id = id

		this.extend(a)

		if err := a.SetStatus(action); err == nil {
			this.renderJson(utils.JsonMessage(true, "", ""))
		} else {
			this.renderJson(utils.JsonMessage(false, "", err.Error()))
		}
	} else {
		this.renderJson(utils.JsonMessage(false, "", "参数错误"))
	}

}

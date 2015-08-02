package controllers

import (
	//"fmt"
	"writer/models"
	"writer/utils"
)

type Article struct {
	Auth
}

func (this *Article) Write() {
	this.Layout = "_frontLayout.html"
	//this.LayoutSections = make(map[string]string)
	this.Data["account"] = this.currentUser

	this.setTplNames("write")
}

func (this *Article) Update() {
	article := new(models.Article)

	article.Id, _ = this.GetInt64("id")
	article.MoreId, _ = this.GetInt64("moreId")
	article.ParentId, _ = this.GetInt64("parentId")
	article.Position, _ = this.GetInt("position")
	article.DocumentId, _ = this.GetInt64("documentId")
	article.Title = this.GetString("title")
	article.Content = this.GetString("content")
	article.Tags = this.GetString("tags")
	//
	if article.Id > 0 {
		this.extend(article)
	} else {
		this.extendEx(article)
	}

	// 提交更新，返回结果
	if err, errs := article.Update(); err == nil {
		this.renderJson(utils.ActionResult(true, article))
	} else {
		this.renderJson(utils.ActionResult(false, errs))
	}

}

/*
方法：读取一条数据
参数：id
*/
func (this *Article) Get() {
	id, err := this.GetInt64("id")

	if err != nil {
		this.renderJson(utils.JsonResult(false, "", "参数错误"))
		return
	}

	//
	a := new(models.Article)
	a.Id = id

	if has, err := a.Get(); err == nil {
		if has {
			this.renderJson(utils.ActionResult(true, a))
		} else {
			this.renderJson(utils.JsonResult(false, "", "数据不存在"))
		}
	} else {
		this.renderJson(utils.JsonResult(false, "", err.Error()))
	}
}

//// @Title List
//// @Description 分页拉取文档列表
//// @Param   size  form  int  false        "每页的记录条数"
//// @Success 200 {object} utils.Response
//// @Failure 200 {object} utils.Response
//// @router /List [post]
//func (this *Article) List() {
//	// 读取分页规则
//	p := new(models.Pagination)

//	if size, err := this.GetInt("size"); err != nil || size == 0 {
//		p.Size = 20
//	}
//	p.Index, _ = this.GetInt("index")
//	// 读取查询条件
//	parentId, _ := this.GetInt64("parentid")
//	tags := this.GetString("tags")

//	// 拉取
//	a := new(models.Article)
//	var as []models.Article
//	var err error

//	// 构造查询字符串
//	cond := "parentId=?"

//	if tags == "" {
//		as, err = a.List(p, cond, parentId)
//	} else {
//		cond += " and tags=?"
//		as, err = a.List(p, cond, parentId, tags)
//	}

//	if err == nil {
//		this.renderJson(utils.JsonResult(true, "", as))
//	} else {
//		this.renderJson(utils.JsonResult(false, "", err.Error()))
//	}
//}

// @Title Position
// @Description 节点位置发生变化
// @Param parentId  form  int  false        "父级id"
// @Param referId  form  int  false        "参考文档id"
func (this *Article) Position() {
	a := new(models.Article)
	a.Id, _ = this.GetInt64("id")
	//
	if a.Id <= 0 {
		this.renderJson(utils.JsonResult(false, "id", "参数错误: id 必须 > 0"))
		return
	}
	// 新的参考点
	a.MoreId, _ = this.GetInt64("moreId")
	// 顺序位置参考点
	a.Position, _ = this.GetInt("referId")

	this.extend(a)

	if ok, err, more := a.SetPosition(); err == nil {
		if ok {
			this.renderJson(utils.JsonResult(true, "", more))
		} else {
			this.renderJson(utils.JsonResult(false, "", ""))
		}
	} else {
		this.renderJson(utils.JsonResult(false, "", err.Error()))
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
		a.ParentId, _ = this.GetInt64("parentId")

		this.extend(a)

		if err := a.SetStatus(action); err == nil {
			this.renderJson(utils.JsonResult(true, "", ""))
		} else {
			this.renderJson(utils.JsonResult(false, "", err.Error()))
		}
	} else {
		this.renderJson(utils.JsonResult(false, "", "参数错误"))
	}

}

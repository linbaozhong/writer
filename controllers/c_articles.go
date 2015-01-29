package controllers

import (
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
	article.Title = this.GetString("title")
	article.Content = this.GetString("content")
	//
	this.extend(article)
	// 提交更新，返回结果
	if err, errs := article.Update(); err == nil {
		this.renderJson(utils.JsonMessage(true, "", ""))
	} else {
		this.renderJson(utils.JsonData(false, "", errs))
	}

}

func (this *Article) Lock() {
	this.status("Articles", "lock")
}

func (this *Article) UnLock() {
	this.status("Articles", "unlock")
}

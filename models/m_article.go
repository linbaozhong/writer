package models

import (
	// "errors"
	// "fmt"
	"github.com/astaxie/beego/validation"
)

type Articles struct {
	Id      int64  `json:"articleId"`
	Title   string `json:"title" valid:"MaxSize(100)"`
	Content string `json:"content"`
	Status  int    `json:"status" valid:"Range(0,1)"`
	Deleted int    `json:"deleted" valid:"Range(0,1)"`
	Creator int64  `json:"creator"`
	Created int64  `json:"created"`
	Updator int64  `json:"updator"`
	Updated int64  `json:"updated"`
	Ip      string `json:"ip" valid:"MaxSize(23)"`
}

// 文章是否存在
func (this *Articles) Exists() (bool, error) {
	return db.Get(this)
}

// 自定义数据验证
func (this *Articles) Valid(v *validation.Validation) {

}

// 新文章
func (this *Articles) Update() (error, []Error) {
	//数据有效性检验
	if d, err := dataCheck(this); err != nil {
		return err, d
	}

	var err error
	if this.Id == 0 {
		_, err = db.Insert(this)
	} else {
		_, err = db.Update(this)
	}
	return err, nil
}

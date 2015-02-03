package models

import (
	// "errors"
	"fmt"
	"github.com/astaxie/beego/validation"
	"strings"
	"zouzhe/utils"
)

type Articles struct {
	Id       int64  `json:"articleId"`
	Title    string `json:"title" valid:"MaxSize(250)"`
	Content  string `json:"content"`
	Tags     string `json:"tags" valid:"MaxSize(250)"`
	ParentId int64  `json:"parentId"`
	Position int    `json:"position"`
	Status   int    `json:"status" valid:"Range(0,1)"`
	Deleted  int    `json:"deleted" valid:"Range(0,1)"`
	Creator  int64  `json:"creator"`
	Created  int64  `json:"created"`
	Updator  int64  `json:"updator"`
	Updated  int64  `json:"updated"`
	Ip       string `json:"ip" valid:"MaxSize(23)"`
}

type Documents struct {
	Id      int64  `json:"articleId"`
	Title   string `json:"title" valid:"MaxSize(250)"`
	Content string `json:"content"`
	Tags    string `json:"tags" valid:"MaxSize(250)"`
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

	session := db.NewSession()
	defer session.Close()
	// add Begin() before any action
	err := session.Begin()

	if err != nil {
		session.Rollback()
		return err, nil
	}

	// 如果是新增，按index腾出位置
	if this.Id == 0 {
		// 找到id=this.Position参考文档的position
		if this.Position > 0 {
			_article := new(Articles)
			if _has, err := session.Id(this.Position).Get(_article); err == nil && _has {
				this.Position = _article.Position
			}
		}
		// 更新其后文档的position
		_, err = session.Exec("update 'articles' set position = position+1 where creator = ? and position > ?", this.Creator, this.Position)
		if err != nil {
			session.Rollback()
			return err, nil
		}
	}

	// 新增 or 修改
	if this.Id == 0 {
		_, err = session.Insert(this)
	} else {
		_, err = session.Id(this.Id).Cols("title", "content", "tags", "updator", "updated", "ip").Update(this)
	}

	if err != nil {
		session.Rollback()
		return err, nil
	}
	// 检查tags是否为空
	if len(this.Tags) != 0 {
		_tags := strings.Split(this.Tags, ",")
		// 从数据库中读取所有标签的id,name
		tags := new(Tags)
		rows, err := session.In("name", _tags).Rows(tags)

		if err != nil {
			session.Rollback()
			return err, nil
		}
		defer rows.Close()

		ids := make([]int64, 0)
		for rows.Next() {
			err = rows.Scan(tags)
			//...
			ids = append(ids, tags.Id)
			// 移除已经存在标签
			if utils.StringsContains(_tags, tags.Name) {
				_tags = utils.RemoveStringSlice(tags.Name, _tags)
			}
		}
		// 遍历标签，找到数据库中不存在的标签，insert into标签表
		for _, v := range _tags {
			_tag := new(Tags)
			_tag.Name = v

			_, err = session.Insert(_tag)
			if err == nil {
				ids = append(ids, _tag.Id)
			} else {
				session.Rollback()
				return err, nil
			}
		}
		// 清除旧的标签-文章的索引
		_del := new(TagArticle)
		_, err = session.Where("articleId = ?", this.Id).Delete(_del)
		if err != nil {
			session.Rollback()
			return err, nil
		}
		// 建立新的标签-文章的索引
		tagArticles := make([]TagArticle, 0)
		for _, id := range ids {
			tagArticles = append(tagArticles, TagArticle{TagId: id, ArticleId: this.Id})
		}
		fmt.Println(tagArticles)
		_, err = session.Insert(tagArticles)
		if err != nil {
			session.Rollback()
			return err, nil
		}
	}
	// add Commit() after all actions
	err = session.Commit()

	return err, nil
}

// 读取
func (this *Articles) Get() (bool, error) {
	return db.Get(this)
}

// 分页列表
func (this *Articles) List(where string, page *Pagination) ([]Articles, error) {
	as := make([]Articles, 0)
	// 符合条件的记录总数
	a := new(Articles)

	if rows, err := db.Where(where).Count(a); err != nil {
		return as, err
	} else {
		getPageCount(rows, page)
	}

	err := db.Where(where).Asc("position").Limit(page.Size, page.Index*page.Size).Find(&as)
	return as, err
}

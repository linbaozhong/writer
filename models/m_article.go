package models

import (
	// "errors"
	"fmt"
	"github.com/astaxie/beego/validation"
	"strings"
	"zouzhe/utils"
)

// article视图
type Article struct {
	Id         int64  `json:"articleId"`
	Title      string `json:"title" valid:"MaxSize(250)"`
	Content    string `json:"content"`
	Tags       string `json:"tags" valid:"MaxSize(250)"`
	ParentId   int64  `json:"parentId"`
	Position   int    `json:"position"`
	DocumentId int64  `json:"documentId"`
	Status     int    `json:"status" valid:"Range(0,1)"`
	Deleted    int    `json:"deleted" valid:"Range(0,1)"`
	Creator    int64  `json:"creator"`
	Created    int64  `json:"created"`
	Updator    int64  `json:"updator"`
	Updated    int64  `json:"updated"`
	Ip         string `json:"ip" valid:"MaxSize(23)"`
}

// Articles表
type Articles struct {
	Id         int64  `json:"articleId"`
	ParentId   int64  `json:"parentId"`
	Position   int    `json:"position"`
	DocumentId int64  `json:"documentId"`
	Status     int    `json:"status" valid:"Range(0,1)"`
	Deleted    int    `json:"deleted" valid:"Range(0,1)"`
	Creator    int64  `json:"creator"`
	Created    int64  `json:"created"`
	Updator    int64  `json:"updator"`
	Updated    int64  `json:"updated"`
	Ip         string `json:"ip" valid:"MaxSize(23)"`
}

// Documents表
type Documents struct {
	Id      int64  `json:"documentId"`
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
func (this *Article) Exists() (bool, error) {
	return db.Get(this)
}

// 自定义数据验证
func (this *Article) Valid(v *validation.Validation) {

}

// 新文章
func (this *Article) Update() (error, []Error) {
	//数据有效性检验
	if d, err := dataCheck(this); err != nil {
		return err, d
	}

	session := db.NewSession()
	defer session.Close()
	// 事务开始
	err := session.Begin()

	if err != nil {
		session.Rollback()
		return err, nil
	}

	// articles对象
	_article := new(Articles)
	_article.Id = this.Id
	_article.ParentId = this.ParentId
	_article.DocumentId = this.DocumentId
	_article.Creator = this.Creator
	_article.Created = this.Created
	_article.Updator = this.Updator
	_article.Updated = this.Updated
	_article.Ip = this.Ip

	// documents对象
	_document := new(Documents)
	_document.Id = this.DocumentId
	_document.Title = this.Title
	_document.Content = this.Content
	_document.Tags = this.Tags
	_document.Creator = this.Creator
	_document.Created = this.Created
	_document.Updator = this.Updator
	_document.Updated = this.Updated
	_document.Ip = this.Ip

	// 如果是新增，按index腾出位置
	if this.Id == 0 {
		// 找到id=this.Position参考文档的position
		var positionSql string
		if this.Position > 0 {
			if _results, err := session.Query("select position from articles where id=?", this.Position); len(_results) > 0 && err == nil {
				_article.Position = utils.Bytes2int(_results[0]["position"])
			}
			positionSql = "update 'articles' set position = position+2 where creator = ? and parentId = ? and position > ?"
		} else {
			positionSql = "update 'articles' set position = position+2 where creator = ? and parentId = ? and position >= ?"
		}
		// 更新其后文档的position
		if _, err = session.Exec(positionSql, this.Updator, this.ParentId, _article.Position); err != nil {
			session.Rollback()
			return err, nil
		}
		// 如果关联文档不存在，创建关联文档
		if this.DocumentId <= 0 {
			// 先 insert documents 附表
			if _, err = session.Insert(_document); err != nil {
				session.Rollback()
				return err, nil
			}
		}
		// insert articles 主表
		_article.DocumentId = _document.Id //主附表映射
		_article.Position += 1
		if _, err = session.Insert(_article); err != nil {
			session.Rollback()
			return err, nil
		}
	} else {
		// Update articles 主表
		if _, err = session.Id(_article.Id).Cols("updator", "updated", "ip").Update(_article); err != nil {
			session.Rollback()
			return err, nil
		}
		// 读取关联的文档id
		if ok, err := this.GetEx(); !ok {
			session.Rollback()
			return err, nil
		}
		// Update documents 附表
		_document.Id = this.DocumentId
		if _, err = session.Id(_document.Id).Cols("title", "content", "tags", "updator", "updated", "ip").Update(_document); err != nil {
			session.Rollback()
			return err, nil
		}
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
	// 提交事务
	err = session.Commit()

	return err, nil
}

// 只读取可见的
func (this *Article) Get() (bool, error) {
	return this._get(true)
}

// 读取
func (this *Article) GetEx() (bool, error) {
	return this._get(false)
}

// 读取
func (this *Article) _get(all bool) (bool, error) {
	// Dal对象
	_dal := &Dal{}
	_dal.Field = "articles.*,documents.title,documents.content"
	_dal.From = "articles,documents"
	_dal.Where = "documents.id = articles.documentId and articles.id=?"

	// 可见的
	if all {
		_dal.Where += fmt.Sprintf(" and articles.status=%d and articles.deleted=%d and documents.status=%d and documents.deleted=%d", Unlock, Undelete, Unlock, Undelete)
	}

	return db.Sql(_dal.Select(), this.Id).Get(this)
}

// 分页列表
func (this *Article) List(page *Pagination, condition string, params ...interface{}) ([]Article, error) {
	return this._list(true, page, condition, params...)
}

// 分页列表
func (this *Article) ListEx(page *Pagination, condition string, params ...interface{}) ([]Article, error) {
	return this._list(false, page, condition, params...)
}

// 分页列表
func (this *Article) _list(all bool, page *Pagination, condition string, params ...interface{}) ([]Article, error) {
	// Dal对象
	_dal := &Dal{}
	_dal.Field = "articles.*,documents.title,documents.content"
	_dal.From = "articles,documents"
	_dal.Where = "documents.id = articles.documentId"
	_dal.OrderBy = "articles.parentId,articles.position"

	// 可见的
	if all {
		_dal.Where += fmt.Sprintf(" and articles.status=%d and articles.deleted=%d and documents.status=%d and documents.deleted=%d", Unlock, Undelete, Unlock, Undelete)
	}
	// 条件
	if strings.TrimSpace(condition) != "" {
		_dal.Where += " and " + condition
	}
	// slice承载返回的结果
	as := make([]Article, 0)
	// 读取符合条件的记录总数
	if rows := _dal.Count(params...); rows > 0 {

		getPageCount(rows, page)

		_dal.Size = page.Size
		_dal.Offset = page.Index * page.Size

		err := db.Sql(_dal.Select(), params...).Find(&as)
		return as, err
	}
	return as, nil
}

// 设置locked、deleted等状态字段的公共方法
func (this *Article) SetStatus(action string) error {
	a := new(Articles)
	a.Id = this.Id
	a.Updated = this.Updated
	a.Updator = this.Updator
	a.Ip = this.Ip

	d := new(Documents)
	d.Updated = this.Updated
	d.Updator = this.Updator
	d.Ip = this.Ip

	switch strings.ToLower(action) {
	case "lock":
		action = "status"
		a.Status = Locked
		d.Status = Locked
	case "unlock":
		action = "status"
		a.Status = Unlock
		d.Status = Unlock
	case "delete":
		action = "deleted"
		a.Deleted = Deleted
		d.Deleted = Deleted
	case "undelete":
		action = "deleted"
		a.Deleted = Undelete
		d.Deleted = Undelete
	}
	// 事务
	session := db.NewSession()
	defer session.Close()
	// 事务开始
	err := session.Begin()

	if err != nil {
		session.Rollback()
		return err
	}
	// 修改articles状态
	_, err = session.Id(a.Id).Cols(action, "updator", "updated", "ip").Update(a)

	if err != nil {
		session.Rollback()
		return err
	}

	// 读取关联的文档id
	if ok, err := this.GetEx(); !ok {
		session.Rollback()
		return err
	}
	d.Id = this.DocumentId

	// 修改documents的状态
	_, err = session.Id(d.Id).Where("creator=?", this.Creator).Cols(action, "updator", "updated", "ip").Update(d)
	if err != nil {
		session.Rollback()
		return err
	}
	// 完成事务提交
	session.Commit()

	return err
}

// 设置文档节点位置
func (this *Article) SetPosition() (bool, error) {
	session := db.NewSession()
	defer session.Close()
	// 事务开始
	err := session.Begin()

	if err != nil {
		session.Rollback()
		return false, err
	}

	// 找到id=this.Position参考文档的position
	var positionSql string
	if this.Position > 0 {
		if _results, err := session.Query("select position from articles where id=?", this.Position); len(_results) > 0 && err == nil {
			this.Position = utils.Bytes2int(_results[0]["position"])
		}
		positionSql = "update 'articles' set position = position+2 where creator = ? and parentId = ? and position > ?"
	} else {
		positionSql = "update 'articles' set position = position+2 where creator = ? and parentId = ? and position >= ?"
	}
	// 更新其后文档的position
	if _, err = session.Exec(positionSql, this.Updator, this.ParentId, this.Position); err != nil {
		session.Rollback()
		return false, err
	}

	a := new(Articles)
	a.Id = this.Id
	a.ParentId = this.ParentId
	a.Position = this.Position + 1

	_, err = session.Id(a.Id).Cols("parentId", "position", "updator", "updated", "ip").Update(a)
	if err != nil {
		session.Rollback()
		return false, err
	}
	// 提交事务
	err = session.Commit()
	return err == nil, err
}

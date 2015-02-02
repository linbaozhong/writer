package models

type Tags struct {
	Id    int64  `json:"tagsId"`
	Name  string `json:"name" valid:"MaxSize(100)"`
	Times int64  `json:"times"`
}
type TagArticle struct {
	TagId     int64
	ArticleId int64
}

// 标签是否存在
func (this *Tags) Exists() (bool, error) {
	return db.Get(this)
}

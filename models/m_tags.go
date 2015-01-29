package models

type Tags struct {
	Id   int64  `json:"tagsId"`
	name string `json:"name" valid:"MaxSize(100)"`
}

// 标签是否存在
func (this *Tags) Exists() (bool, error) {
	return db.Get(this)
}

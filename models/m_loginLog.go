package models

// 登录日志
type LoginLog struct {
	Id        int64  `json:"loginId"`
	AccountId int64  `json:"accoundId"`
	Updated   int64  `json:"loginTime"`
	Ip        string `json:"ip" valid:"MaxSize(23)"`
}

//
func (this *LoginLog) Post() (int64, error) {
	return db.Insert(this)
}

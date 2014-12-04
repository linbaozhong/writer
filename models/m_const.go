package models

const (
	Role_System        = 0 //超级管理员
	Role_Administrator = 1 //管理员
	Role_Founder       = 2 //创建者
	Role_Author        = 3 //作者
	Role_Reader        = 4 //读者
	Role_Guest         = 5 //游客

	Open_Alipay = 0 //支付宝用户
	Open_QQ     = 1 //qq用户
	Open_Weibo  = 2 //微博用户

	Private  = 0 //私有的
	Internal = 1 //内部的
	Public   = 2 //公开的
	Free     = 3 //自由的

	Undelete = 0 //未删除
	Deleted  = 1 //删除

	Unlock = 0 //未锁定
	Locked = 1 //锁定
)

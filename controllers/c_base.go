package controllers

import (
	"github.com/astaxie/beego"
)

type Base struct {
	beego.Controller
}

//输出字符串
func (b *Base) toString(arg string) {
	b.Ctx.Output.Body([]byte(arg))
}

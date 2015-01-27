package controllers

import (
//"html/template"
)

type Home struct {
	Front
}

func (this *Home) Get() {
	this.setTplNames("index")
}

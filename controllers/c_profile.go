package controllers

type Profile struct {
	Auth
}

func (this *Profile) Get() {
	//this.Data["Title"] = appTitle + " - 线路设计"
	this.TplNames = "profile/line.html"
}

func (this *Profile) Line() {
	//this.Data["Title"] = appTitle + " - 线路设计"
	this.TplNames = "line.html"
}

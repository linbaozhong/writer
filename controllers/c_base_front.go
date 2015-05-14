package controllers

type Front struct {
	Base
}

func (this *Front) Prepare() {
	this.Base.Prepare()
	this.allowRequest()

	this.Layout = "_frontLayout.html"
	this.LayoutSections = make(map[string]string)
	// this.LayoutSections["Head"] = "_head.html"
	// this.LayoutSections["Header"] = "_indexHeader.html"
	// this.LayoutSections["Login"] = "_login.html"
	// this.LayoutSections["Footer"] = "_footer.html"
	// this.LayoutSections["Scripts"] = "_scripts.html"

	this.Data["account"] = this.currentUser
}

func (this *Front) Finish() {

}

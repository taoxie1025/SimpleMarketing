package controllers

import (
	"email_action/store"
	"github.com/astaxie/beego"
)

type HomePageController struct {
	beego.Controller
	Store *store.Store
}

func (c *HomePageController) Get() {
	c.TplName = "index.tpl"
}

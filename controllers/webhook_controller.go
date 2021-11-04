package controllers

import (
	"email_action/store"
	"github.com/astaxie/beego"
	"io/ioutil"
)

type WebhookController struct {
	beego.Controller
	Store *store.Store
}

func (c *WebhookController) StripeEventListener() {
	log.Infof("StripeEventListener():")
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		log.Errorf("Error reading request body: %v\n", err)
		c.Abort("400")
	}
	header := c.Ctx.Input.Header("Stripe-Signature")
	if err := c.Store.ProcessStripeEvents(body, header); err != nil {
		c.Abort("401")
		return
	}
	c.ServeJSON()
}

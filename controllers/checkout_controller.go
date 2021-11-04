package controllers

import (
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type CheckoutController struct {
	beego.Controller
	Store *store.Store
}

func (c *CheckoutController) CreateSession() {
	var req models.NewCheckoutSession
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse new NewCheckoutSession"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}

	sess, err := c.Store.CreateCheckoutSession(req.StripeCustomerId, req.PriceId, req.Plan, req.Coupon, req.Email, req.SuccessUrl, req.CancelUrl)
	if err != nil {
		errMsg := "failed to create new user"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}

	c.Data["json"] = sess
	c.ServeJSON()
}

func (c *CheckoutController) CancelSubscription() {
	email := c.Ctx.Input.Param(":email")
	subscriptionId := c.Ctx.Input.Param(":subscriptionId")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	if requesterEmail != email {
		c.Abort("400")
	}

	err := c.Store.GetStripeAdapter().CancelSubscription(subscriptionId)
	if err != nil {
		c.Abort("500")
	}
	err = c.Store.UpdateSubscriptionPlan(email, 0, 0, "", "", false)
	if err != nil {
		c.Abort("500")
	}
	c.ServeJSON()
}

func (c *CheckoutController) UpdateSubscription() {
	email := c.Ctx.Input.Param(":email")
	subscriptionId := c.Ctx.Input.Param(":subscriptionId")

	var req models.UpdateSubscriptionRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse UpdateSubscriptionRequest"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}

	err = c.Store.GetStripeAdapter().UpdateSubscription(subscriptionId, req.PriceId)
	if err != nil {
		c.Abort("500")
	}

	targetPlan := models.SubscriptionPlanFree
	if req.PriceId == beego.AppConfig.String("pro_plan_price_id") {
		targetPlan = models.SubscriptionPlanPro
	} else if req.PriceId == beego.AppConfig.String("ultra_plan_price_id") {
		targetPlan = models.SubscriptionPlanUltra
	}
	err = c.Store.UpdateSubscriptionPlan(email, targetPlan, 0, subscriptionId, req.PriceId, false)
	if err != nil {
		c.Abort("500")
	}
	c.ServeJSON()
}

func (c *CheckoutController) ReadCoupon() {
	couponId := c.Ctx.Input.Query("couponId")
	coupon, err := c.Store.ReadCoupon(couponId)
	if err != nil || coupon == nil {
		c.Abort("404")
	}
	c.Data["json"] = coupon
	c.ServeJSON()
}

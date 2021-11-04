package controllers

import (
	"email_action/mail"
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type EmailController struct {
	beego.Controller
	Store  *store.Store
	Mailer *mail.SesAdapter
}

func (c *EmailController) IsEmailVerified() {
	email := c.Ctx.Input.Param(":email")
	if email == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	log.Infof("IsEmailVerified(): %s", email)

	if isVerified, err := c.Mailer.IsEmailVerified(email); err == nil {
		if isVerified {
			c.Data["json"] = &models.EmailVerificationResult{
				Email:      email,
				IsVerified: true,
			}
			c.ServeJSON()
		} else {
			c.Data["json"] = &models.EmailVerificationResult{
				Email:      email,
				IsVerified: false,
			}
			c.ServeJSON()
		}
	}
	c.Abort("500")
}

func (c *EmailController) SendVerificationEmail() {
	email := c.Ctx.Input.Param(":email")
	if email == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	log.Infof("SendVerificationEmail(): %s", email)

	if err := c.Mailer.VerifyEmail(email); err != nil {
		c.Abort("500")
	}
	c.ServeJSON()
}

func (c *EmailController) ContactUs() {
	log.Infof("ContactUs():")
	var req models.ContactUsRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse ContactUsRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}

	contactEmail := beego.AppConfig.String("contact_us_email")
	if err := c.Mailer.SendEmail(req.Email, contactEmail, req.Subject, req.Message, req.Message); err != nil {
		c.Abort("500")
	}
	c.ServeJSON()
}

func (c *EmailController) SendEmail() {
	log.Infof("SendEmail():")
	var req models.SendEmailRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse SendEmailRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	// increment quota for testing purpose as well
	c.Store.AddEmailUsageInCycle(req.Email, 1)
	message := req.Message + c.Store.GetUnsubscribeLink(req.ProjectId, req.ProjectName)
	if err := c.Mailer.SendEmail(req.Email, req.To, req.Subject, message, message); err != nil {
		c.Abort("500")
	}
	c.ServeJSON()
}

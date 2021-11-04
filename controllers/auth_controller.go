package controllers

import (
	"email_action/mail"
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type AuthController struct {
	beego.Controller
	Store  *store.Store
	Mailer *mail.SesAdapter
}

func (c *AuthController) AuthUser() {
	var req models.AuthRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse auth data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}

	claims, token, err := c.Store.AuthUser(&req)
	if err != nil {
		c.Abort("400")
	}
	resp := &models.AuthResponse{
		Token:  token,
		Claims: claims,
	}

	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *AuthController) ResetPassword() {
	var req models.ResetPasswordRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse ResetPasswordRequest"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}

	resetLink, err := c.Store.ResetPassword(req.Email)
	if err != nil {
		c.Abort("400")
	}
	from := beego.AppConfig.String("no_reply_email")
	subject := "Reset the Password for Your Simple Marketing Account"
	htmlBody := "Please click the link to reset your password:\n" + resetLink
	c.Mailer.SendEmail(from, req.Email, subject, htmlBody, htmlBody)
	c.ServeJSON()
}

func (c *AuthController) IsTokenValid() {
	var req models.ValidateTokenRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse ValidateTokenRequest"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	err = c.Store.ReadPasswordReset(req.Token)
	if err != nil {
		c.Abort("400")
	}
	c.ServeJSON()
}

func (c *AuthController) UpdatePassword() {
	var req models.UpdatePasswordTokenRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse UpdatePasswordTokenRequest"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	err = c.Store.UpdatePassword(req)
	if err != nil {
		c.Abort("400")
	}
	c.ServeJSON()
}

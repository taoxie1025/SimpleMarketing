package controllers

import (
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
	Store *store.Store
}

func (c *UserController) CreateUser() {
	var req models.NewUserRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse new user data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	if req.Password == "" {
		errMsg := "missing password"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	ip := c.Ctx.Request.Host
	user, err := c.Store.CreateUser(req.Email, req.Password, req.FirstName, req.LastName, ip, req.Address, req.PhoneNumber)
	if err != nil {
		errMsg := "failed to create new user"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}

	authReq := &models.AuthRequest{
		Email:    user.Email,
		Password: req.Password,
	}
	claims, token, err := c.Store.AuthUser(authReq)
	if err != nil {
		c.Abort("400")
	}
	resp := models.AuthResponse{
		Token:  token,
		Claims: claims,
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *UserController) UpdateUserBasicInfo() {
	var req models.UpdateUserRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse new user data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	if req.Email != req.UserInfo.Email {
		log.Errorf("UpdateUserBasicInfo(): %v", unauthorizedErr)
		c.Abort("401")
	}
	user, err := c.Store.UpdateUserBasicInfo(&req.UserInfo)
	if err != nil {
		log.Errorf("UpdateUserBasicInfo(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = user
	c.ServeJSON()
}

func (c *UserController) ReadUser() {
	log.Infof("ReadUser():")
	email := c.Ctx.Input.Param(":email")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	if email == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	if email != requesterEmail {
		log.Errorf("ReadUser(): %v", unauthorizedErr)
		c.Abort("401")
	}
	user, err := c.Store.ReadUser(email)
	if err != nil {
		log.Errorf("ReadUser(): error = %+v", err)
		c.Abort("500")
	}
	user.PasswordHash = ""
	c.Data["json"] = user
	c.ServeJSON()
}

func (c *UserController) GetAccountInfo() {
	log.Infof("GetAccountInfo():")
	email := c.Ctx.Input.Param(":email")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	if email == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	if email != requesterEmail {
		log.Errorf("ReadUser(): %v", unauthorizedErr)
		c.Abort("401")
	}
	user, err := c.Store.ReadUser(email)
	accountInfo := models.NewUserAccountInfo(user)
	if err != nil {
		log.Errorf("GetAccountInfo(): error = %+v", err)
		c.Abort("500")
	}
	user.PasswordHash = ""
	c.Data["json"] = accountInfo
	c.ServeJSON()
}

func (c *UserController) UpdatePassword() {
	var req models.UpdatePasswordRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse UpdatePasswordRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	err = c.Store.ChangePassword(req)
	if err != nil {
		log.Errorf("UpdatePassword(): error = %+v", err)
		c.Abort("500")
	}
	c.ServeJSON()
}

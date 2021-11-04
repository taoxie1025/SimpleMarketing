package controllers

import (
	"email_action/mail"
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type AdminController struct {
	beego.Controller
	Store  *store.Store
	Mailer *mail.SesAdapter
}

func (c *AdminController) AdminReadTickets() {
	log.Infof("AdminReadTickets():")
	var req models.AdminReadTicketsRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse AdminReadTicketsRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	readTicketsResponse, err := c.Store.AdminReadTickets(req.Email, req.Token, req.PageSize)
	if err != nil {
		log.Errorf("AdminReadTickets(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = readTicketsResponse
	c.ServeJSON()
}

func (c *AdminController) AdminReadUsers() {
	log.Infof("AdminReadUsers():")
	var req models.AdminReadUsersRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse AdminReadUsersRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	readUsersResponse, err := c.Store.AdminReadUsers(req.Email, req.Token, req.PageSize)
	if err != nil {
		log.Errorf("AdminReadUsers(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = readUsersResponse
	c.ServeJSON()
}

func (c *AdminController) AdminSearchUsers() {
	log.Infof("AdminSearchUsers():")
	emailFilter := c.Ctx.Input.Query("emailFilter")
	if emailFilter == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}

	resp, err := c.Store.SearchUsers(emailFilter)
	if err != nil {
		log.Errorf("AdminSearchUsers(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *AdminController) AdminSearchTickets() {
	log.Infof("AdminSearchTickets():")
	ticketFilter := c.Ctx.Input.Query("ticketFilter")
	if ticketFilter == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}

	resp, err := c.Store.SearchTickets(ticketFilter)
	if err != nil {
		log.Errorf("AdminSearchTickets(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *AdminController) AdminUpdateUser() {
	log.Infof("AdminUpdateUser():")
	var req models.AdminUpdateUserRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse AdminUpdateUserRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	resp, err := c.Store.AdminUpdateUser(&req)
	if err != nil {
		log.Errorf("AdminUpdateUser(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *AdminController) AdminReadSubscribers() {
	log.Infof("AdminReadSubscribers():")
	var req models.AdminReadSubscribersRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse AdminReadSubscribersRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	readSubscribersResponse, err := c.Store.AdminReadSubscribers(req.Email, req.Token, req.PageSize)
	if err != nil {
		log.Errorf("AdminReadSubscribers(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = readSubscribersResponse
	c.ServeJSON()
}

func (c *AdminController) AdminSearchSubscribers() {
	log.Infof("AdminSearchSubscribers():")
	subscriberFilter := c.Ctx.Input.Query("subscriberFilter")
	if subscriberFilter == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}

	subscribers, err := c.Store.AdminSearchSubscribers(subscriberFilter)
	if err != nil {
		log.Errorf("AdminSearchSubscribers(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = subscribers
	c.ServeJSON()
}

func (c *AdminController) AdminUpdateSubscriber() {
	log.Infof("AdminUpdateSubscriber():")
	var req models.AdminUpdateSubscriberRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse AdminUpdateSubscriberRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	resp, err := c.Store.AdminUpdateSubscriber(&req)
	if err != nil {
		log.Errorf("AdminUpdateSubscriber(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

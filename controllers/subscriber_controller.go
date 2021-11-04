package controllers

import (
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type SubscriberController struct {
	beego.Controller
	Store *store.Store
}

func (c *SubscriberController) CreateSubscriber() {
	log.Infof("CreateSubscriber():")
	var req models.NewSubscriberRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse NewSubscriberRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	if req.Email == "" || req.ProjectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	subscriber, err := c.Store.CreateSubscriber(req.ProjectId, req.Email, req.FirstName, req.LastName)
	if err != nil {
		errMsg := "failed to create new subscriber"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}
	c.Data["json"] = subscriber
	c.ServeJSON()
}

func (c *SubscriberController) UpdateSubscriber() {
	log.Infof("UpdateSubscriber():")
	var req models.UpdateSubscriberRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse UpdateSubscriberRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	if req.Email == "" || req.ProjectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	subscriber, err := c.Store.UpdateSubscriber(&req)
	if err != nil {
		errMsg := "failed to create update subscriber"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}
	c.Data["json"] = subscriber
	c.ServeJSON()
}

func (c *SubscriberController) DeleteSubscriber() {
	log.Infof("DeleteSubscriber():")
	projectId := c.Ctx.Input.Param(":projectId")
	email := c.Ctx.Input.Param(":email")
	if email == "" || projectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	err := c.Store.DeleteSubscriber(projectId, email)
	if err != nil {
		log.Errorf("DeleteSubscriber(): error = %+v", err)
		c.Abort("500")
	}
	c.ServeJSON()
}

func (c *SubscriberController) ReadSubscribers() {
	log.Infof("ReadSubscribers():")
	var req models.ReadSubscribersRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse NewSubscriberRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	if req.Email == "" || req.ProjectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	if req.PageSize <= 0 {
		defaultPageSize := beego.AppConfig.DefaultInt("default_read_subscribers_page_size", 30)
		req.PageSize = defaultPageSize
	}
	resp, err := c.Store.ReadSubscribers(req.ProjectId, req.Token, req.PageSize)
	if err != nil {
		log.Errorf("ReadSubscribers(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *SubscriberController) SearchSubscribers() {
	log.Infof("SearchSubscribers():")
	projectId := c.Ctx.Input.Query("projectId")
	emailFilter := c.Ctx.Input.Query("emailFilter")
	if emailFilter == "" || projectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}

	resp, err := c.Store.SearchSubscribers(projectId, emailFilter)
	if err != nil {
		log.Errorf("SearchSubscribers(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

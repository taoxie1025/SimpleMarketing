package controllers

import (
	"email_action/mail"
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type TicketController struct {
	beego.Controller
	Store  *store.Store
	Mailer *mail.SesAdapter
}

func (c *TicketController) CreateTicket() {
	log.Infof("CreateTicket():")
	var req models.NewTicketRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse NewTicketRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	if req.Email == "" || req.Title == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	ticket, err := c.Store.CreateTicket(req.Email, req.ProjectId, req.ProjectName, req.Name, req.Title, req.Body, req.TicketType)
	if err != nil {
		errMsg := "failed to create new ticket"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}
	c.Data["json"] = ticket
	c.ServeJSON()
}

func (c *TicketController) UpdateTicket() {
	log.Infof("UpdateTicket():")
	var req models.UpdateTicketRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse UpdateTicketRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	if req.Email == "" || req.Ticket == nil || req.Ticket.TicketId == "" || req.Ticket.Email == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	ticket, err := c.Store.UpdateTicket(req.Ticket)
	if err != nil {
		errMsg := "failed to update ticket"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}
	c.Data["json"] = ticket
	c.ServeJSON()
}

func (c *TicketController) ReadTickets() {
	log.Infof("ReadTickets():")
	var req models.ReadTicketsRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse ReadTicketsRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	readTicketsResponse, err := c.Store.ReadTickets(req.Email, req.Token, req.PageSize)
	if err != nil {
		log.Errorf("ReadTickets(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = readTicketsResponse
	c.ServeJSON()
}

func (c *TicketController) DeleteTicket() {
	log.Infof("DeleteTicket():")
	email := c.Ctx.Input.Param(":email")
	ticketId := c.Ctx.Input.Param(":ticketId")
	if email == "" || ticketId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	err := c.Store.DeleteTicket(email, ticketId)
	if err != nil {
		log.Errorf("DeleteTicket(): error = %+v", err)
		c.Abort("500")
	}
	c.ServeJSON()
}

func (c *TicketController) CreateComment() {
	log.Infof("CreateComment():")
	var req models.NewCommentRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse NewCommentRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	if req.Email == "" || req.Body == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	comment := models.NewComment(req.Email, req.TicketId, req.Name, req.Body)
	comment, err = c.Store.CreateComment(req.Email, req.TicketId, comment)
	if err != nil {
		errMsg := "failed to create new comment"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}
	c.Data["json"] = comment
	c.ServeJSON()
}

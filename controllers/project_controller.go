package controllers

import (
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
	"math"
	"time"
)

type ProjectController struct {
	beego.Controller
	Store *store.Store
}

func (c *ProjectController) CreateProject() {
	log.Infof("CreateProject():")
	var req models.NewProjectRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse NewProjectRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Data["content"] = errMsg
		c.Abort("500")
	}
	if req.Email == "" || req.Name == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Data["content"] = errMsg
		c.Abort("400")
	}
	req.Interval = int64(math.Max(float64(req.Interval), float64((time.Hour * 24).Milliseconds())))
	project, err := c.Store.CreateProject(req.Email, req.Name, req.Intro, req.OutgoingEmail, req.AvatarUrl, req.BackgroundImageUrl, req.Author, req.Interval, req.SubscriptionType)
	if err != nil {
		errMsg := "failed to create new project"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}
	c.Data["json"] = project
	c.ServeJSON()
}

func (c *ProjectController) UpdateProject() {
	log.Infof("UpdateProject():")
	var req models.UpdateProjectRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse UpdateProjectRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	if req.Email != req.Project.Email {
		log.Errorf("UpdateProject(): %v", unauthorizedErr)
		c.Abort("401")
	}
	req.Project.Interval = int64(math.Max(float64(req.Project.Interval), float64((time.Hour * 24).Milliseconds())))
	project, err := c.Store.SaveProject(&req.Project)
	if err != nil {
		log.Errorf("UpdateProject(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = project
	c.ServeJSON()
}

func (c *ProjectController) ReadProject() {
	log.Infof("ReadProject(): %s", c.Ctx.Input.RequestBody)
	email := c.Ctx.Input.Param(":email")
	projectId := c.Ctx.Input.Param(":projectId")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	if email == "" || projectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	if email != requesterEmail {
		log.Errorf("ReadProject(): %v", unauthorizedErr)
		c.Abort("401")
	}
	project, err := c.Store.ReadProject(email, projectId)
	if err != nil {
		log.Errorf("ReadProject(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = project
	c.ServeJSON()
}

func (c *ProjectController) DeleteProject() {
	log.Infof("DeleteProject():")
	email := c.Ctx.Input.Param(":email")
	projectId := c.Ctx.Input.Param(":projectId")
	if email == "" || projectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	err := c.Store.DeleteProject(email, projectId)
	if err != nil {
		log.Errorf("DeleteProject(): error = %+v", err)
		c.Abort("500")
	}
	c.ServeJSON()
}

func (c *ProjectController) ReadProjects() {
	log.Infof("ReadProjects(): %s", c.Ctx.Input.RequestBody)
	var req models.BatchReadProjectRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse BatchReadProjectRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	projects, err := c.Store.ReadProjects(req.Email, req.ProjectIds)
	if err != nil {
		log.Errorf("ReadProjects(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = projects
	c.ServeJSON()
}

func (c *ProjectController) ReadProjectBrief() {
	log.Infof("ReadProjectBrief(): %s", c.Ctx.Input.RequestBody)
	projectId := c.Ctx.Input.Param(":projectId")

	projectBrief, err := c.Store.ReadProjectBrief(projectId)
	if err != nil {
		log.Errorf("ReadProjectBrief(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = projectBrief
	c.ServeJSON()
}

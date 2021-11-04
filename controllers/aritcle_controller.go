package controllers

import (
	"email_action/models"
	"email_action/store"
	"encoding/json"
	"github.com/astaxie/beego"
)

type ArticleController struct {
	beego.Controller
	Store *store.Store
}

func (c *ArticleController) CreateArticle() {
	log.Infof("CreateArticle():")
	var req models.NewArticleRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse NewArticleRequest data"
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
	article, err := c.Store.CreateArticle(req.Email, req.ProjectId, req.Title, req.HtmlBody, req.TextBody)
	if err != nil {
		errMsg := "failed to create new article"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("400")
	}
	c.Data["json"] = article
	c.ServeJSON()
}

func (c *ArticleController) UpdateArticle() {
	log.Infof("UpdateArticle():")
	var req models.UpdateArticleRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse UpdateProjectRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	if req.Email != req.Article.Email {
		log.Errorf("UpdateArticle(): %v", unauthorizedErr)
		c.Abort("401")
	}
	article, err := c.Store.UpdateArticle(&req.Article)
	if err != nil {
		log.Errorf("UpdateArticle(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = article
	c.ServeJSON()
}

func (c *ArticleController) ReadArticle() {
	log.Infof("ReadArticle(): %s", c.Ctx.Input.RequestBody)
	email := c.Ctx.Input.Param(":email")
	articleId := c.Ctx.Input.Param(":articleId")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	if email == "" || articleId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	if email != requesterEmail {
		log.Errorf("UpdateArticle(): %v", unauthorizedErr)
		c.Abort("401")
	}
	article, err := c.Store.ReadArticle(email, articleId)
	if err != nil {
		log.Errorf("ReadArticle(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = article
	c.ServeJSON()
}

func (c *ArticleController) ReadArticles() {
	log.Infof("ReadArticles(): %s", c.Ctx.Input.RequestBody)
	var req models.BatchReadArticleRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		errMsg := "failed to parse BatchReadArticleRequest data"
		log.Errorf("%s, error = %+v", errMsg, err)
		c.Abort("500")
	}
	articles, err := c.Store.ReadArticles(req.Email, req.ArticleIds)
	if err != nil {
		log.Errorf("ReadArticles(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = articles
	c.ServeJSON()
}

func (c *ArticleController) DeleteArticle() {
	log.Infof("DeleteArticle():")
	email := c.Ctx.Input.Param(":email")
	articleId := c.Ctx.Input.Param(":articleId")
	projectId := c.Ctx.Input.Param(":projectId")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")

	if email == "" || articleId == "" || projectId == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	if email != requesterEmail {
		log.Errorf("DeleteArticle(): %v", unauthorizedErr)
		c.Abort("401")
	}

	err := c.Store.DeleteArticle(email, projectId, articleId)
	if err != nil {
		log.Errorf("ReadArticle(): error = %+v", err)
		c.Abort("500")
	}
	c.ServeJSON()
}

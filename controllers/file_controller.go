package controllers

import (
	"email_action/models"
	"email_action/store"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type FileController struct {
	beego.Controller
	Store *store.Store
}

func (c *FileController) FileUpload() {
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	if requesterEmail == "" {
		errMsg := "invalid username"
		log.Errorf("FileUpload(): %s", errMsg)
		c.Abort("400")
	}

	mFile, header, err := c.GetFile("file")
	if err != nil {
		log.Errorf("FileUpload(): error = %+v", err)
		c.Abort("400")
	}
	if mFile == nil || header.Filename == "" {
		errMsg := "invalid file or file name"
		log.Errorf("FileUpload(): %s", errMsg)
		c.Abort("400")
	}

	uploadDir := "../upload/" + time.Now().Format("2006/01/16/")
	err = os.MkdirAll(uploadDir, 777)
	if err != nil {
		log.Errorf("FileUpload() error = %v", err)
		c.Abort("500")
	}

	fileName := strings.ReplaceAll(header.Filename, " ", "")
	path := uploadDir + fileName
	defer mFile.Close()
	err = c.SaveToFile("file", path)
	if err != nil {
		log.Errorf("FileUpload(): error = %v", err)
		c.Abort("500")
	}
	bytes, err := getByteArray(path)
	if err != nil {
		c.Abort("500")
	}

	fileSizeLimit := beego.AppConfig.DefaultInt("filesizelimit", 3000000)
	if len(bytes) > fileSizeLimit {
		log.Errorf("FileUpload(): file size exceeds %d/%d", len(bytes), fileSizeLimit)
		c.Abort("400")
	}

	url, err := c.Store.UploadFile(bytes, requesterEmail, fileName)
	if err != nil {
		log.Errorf("FileUpload(): error = %+v", err)
		c.Abort("500")
	}
	file := &models.File{
		Url:          url,
		ThumbnailUrl: url + "-thumbnail",
	}
	c.Data["json"] = file
	c.ServeJSON()
}

func (c *FileController) ListFiles() {
	log.Infof("ListFiles():")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	email := c.Ctx.Input.Param(":email")
	if email == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	if email != requesterEmail {
		log.Errorf("ListFiles(): %v", unauthorizedErr)
		c.Abort("401")
	}
	files, err := c.Store.ListFiles(email)
	if err != nil {
		log.Errorf("ListFiles(): error = %+v", err)
		c.Abort("500")
	}
	c.Data["json"] = files
	c.ServeJSON()
}

func (c *FileController) DeleteFile() {
	log.Infof("DeleteFile():")
	requesterEmail := c.Ctx.Input.Query("requesterEmail")
	email := c.Ctx.Input.Param(":email")
	fileName := c.Ctx.Input.Param(":fileName")

	if email == "" {
		errMsg := "invalid request"
		log.Errorf("%s", errMsg)
		c.Abort("400")
	}
	if email != requesterEmail {
		log.Errorf("DeleteFile(): %v", unauthorizedErr)
		c.Abort("401")
	}
	err := c.Store.DeleteFile(email, fileName)
	if err != nil {
		log.Errorf("DeleteFile(): error = %+v", err)
		c.Abort("500")
	}
	c.ServeJSON()
}

func getByteArray(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("getByteArray() error = %v", err)
		return nil, err
	}
	return data, nil
}

package models

import "time"

type Article struct {
	Email     string `json:"email"`     // the user email it belongs to, primary key
	ArticleId string `json:"articleId"` // ArticleId is the sort key
	ProjectId string `json:"projectId"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	IsLive    bool   `json:"isLive"`
	Title     string `json:"title"`
	HtmlBody  string `json:"htmlBody"`
	TextBody  string `json:"textBody"`
}

func NewArticle(email, projectId, title, htmlBody, textBody string) *Article {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return &Article{
		Email:     email,
		ProjectId: projectId,
		ArticleId: GenerateArticleUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		Title:     title,
		HtmlBody:  htmlBody,
		TextBody:  textBody,
	}
}

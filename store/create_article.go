package store

import (
	"email_action/models"
)

func (s *Store) CreateArticle(email, projectId, title, htmlBody, textBody string) (*models.Article, error) {
	log.Infof("CreateArticle(): %s, %s", email, projectId)
	article := models.NewArticle(email, projectId, title, htmlBody, textBody)
	err := s.dynamodbAdapter.SaveArticle(article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

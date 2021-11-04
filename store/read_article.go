package store

import (
	"email_action/models"
)

func (s *Store) ReadArticle(email, articleId string) (*models.Article, error) {
	log.Infof("ReadArticle(): %s, %s, %s", email, articleId)
	article, err := s.dynamodbAdapter.ReadArticle(email, articleId)
	if err != nil {
		return nil, err
	}
	return article, nil
}

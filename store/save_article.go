package store

import (
	"email_action/models"
)

func (s *Store) SaveArticle(article *models.Article) (*models.Article, error) {
	log.Infof("SaveArticle(): %v", article)
	err := s.dynamodbAdapter.SaveArticle(article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

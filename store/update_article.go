package store

import (
	"email_action/models"
)

func (s *Store) UpdateArticle(article *models.Article) (*models.Article, error) {
	log.Infof("UpdateArticle(): %v", article)
	err := s.dynamodbAdapter.UpdateArticle(article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

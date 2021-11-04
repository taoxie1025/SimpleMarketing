package store

import (
	"email_action/models"
)

func (s *Store) CreateComment(email, ticketId string, comment *models.Comment) (*models.Comment, error) {
	log.Infof("CreateComment(): %s, %s", email, ticketId)
	err := s.dynamodbAdapter.AppendComment(email, ticketId, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

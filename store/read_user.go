package store

import (
	"email_action/models"
)

func (s *Store) ReadUser(email string) (*models.User, error) {
	log.Infof("ReadUser(): %s", email)
	user, err := s.dynamodbAdapter.ReadUser(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

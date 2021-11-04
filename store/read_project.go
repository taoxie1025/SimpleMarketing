package store

import (
	"email_action/models"
)

func (s *Store) ReadProject(email, projectId string) (*models.Project, error) {
	log.Infof("ReadProject(): %s, %s", email, projectId)
	project, err := s.dynamodbAdapter.ReadProject(email, projectId)
	if err != nil {
		return nil, err
	}
	return project, nil
}

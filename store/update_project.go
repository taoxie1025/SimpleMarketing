package store

import (
	"email_action/models"
	"time"
)

func (s *Store) SaveProject(project *models.Project) (*models.Project, error) {
	log.Infof("SaveProject(): %v", project)
	project.UpdatedAt = time.Now().Unix() * 1000
	err := s.dynamodbAdapter.SaveProject(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

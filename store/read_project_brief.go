package store

import (
	"email_action/models"
)

func (s *Store) ReadProjectBrief(projectId string) (*models.ProjectBrief, error) {
	log.Infof("ReadProjectBrief(): %s", projectId)
	projectBrief, err := s.dynamodbAdapter.ReadProjectBrief(projectId)
	if err != nil {
		return nil, err
	}
	return projectBrief, nil
}

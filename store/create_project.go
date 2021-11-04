package store

import (
	"email_action/models"
)

func (s *Store) CreateProject(email, name, intro, outgoingEmail, avatarUrl, backgroundImageUrl, author string, interval int64, subscriptionType int) (*models.Project, error) {
	log.Infof("CreateProject(): %s, %s, %d", email, name, subscriptionType)
	project := models.NewProject(email, name, intro, outgoingEmail, avatarUrl, backgroundImageUrl, author, interval, subscriptionType)
	err := s.dynamodbAdapter.CreateProjectIfNotExist(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

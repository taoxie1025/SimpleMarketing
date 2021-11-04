package store

import "email_action/models"

func (s *Store) SearchSubscribers(projectId, emailFilter string) ([]*models.Subscriber, error) {
	log.Infof("SearchSubscribers(): %s, %s", projectId, emailFilter)
	subscribers, err := s.dynamodbAdapter.SearchSubscribers(projectId, emailFilter)
	if err != nil {
		return nil, err
	}
	return subscribers, nil
}

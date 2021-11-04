package store

import (
	"email_action/models"
)

func (s *Store) UpdateSubscriber(req *models.UpdateSubscriberRequest) (*models.Subscriber, error) {
	log.Infof("UpdateSubscriber(): %s, %s", req.Email, req.ProjectId)
	err := s.dynamodbAdapter.UpdateSubscriberBasic(req)
	if err != nil {
		return nil, err
	}
	subscriber, err := s.dynamodbAdapter.ReadSubscriber(req.Email, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}

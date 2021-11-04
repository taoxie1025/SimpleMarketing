package store

import "email_action/models"

func (s *Store) ReadSubscribers(projectId, token string, pageSize int) (*models.ReadSubscribersResponse, error) {
	log.Infof("ReadSubscribers(): %s, %s, %d", projectId, token, pageSize)
	subscribers, nextToken, err := s.dynamodbAdapter.ReadSubscribers(projectId, token, pageSize)
	if err != nil {
		return nil, err
	}
	readSubscriberResp := &models.ReadSubscribersResponse{
		Token:       nextToken,
		Subscribers: subscribers,
	}
	return readSubscriberResp, nil
}

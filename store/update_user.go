package store

import (
	"email_action/models"
)

func (s *Store) UpdateUserBasicInfo(userInfo *models.UserInfo) (*models.UserInfo, error) {
	log.Infof("UpdateUserBasicInfo(): %v", userInfo)
	err := s.dynamodbAdapter.UpdateUserBasicInfo(userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

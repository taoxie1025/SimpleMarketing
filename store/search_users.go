package store

import "email_action/models"

func (s *Store) SearchUsers(emailFilter string) ([]*models.UserAccountInfo, error) {
	log.Infof("SearchUsers(): %s", emailFilter)
	users, err := s.dynamodbAdapter.SearchUsers(emailFilter)
	userAccountInfos := []*models.UserAccountInfo{}
	for _, user := range users {
		userAccountInfo := models.NewUserAccountInfo(user)
		userAccountInfos = append(userAccountInfos, userAccountInfo)
	}
	if err != nil {
		return nil, err
	}
	return userAccountInfos, nil
}

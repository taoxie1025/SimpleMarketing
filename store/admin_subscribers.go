package store

import (
	"email_action/models"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

func (s *Store) AdminReadSubscribers(adminEmail, token string, pageSize int) (
	*models.AdminReadSubscribersResponse, error) {
	log.Infof("AdminReadSubscribers(): %s, %s, %d", adminEmail, token, pageSize)
	if !s.hasAdminPrivilege(adminEmail) {
		return nil, errors.New("no admin privilege")
	}
	return s.adminReadSubscribers(token, pageSize)
}

func (s *Store) AdminUpdateSubscriber(request *models.AdminUpdateSubscriberRequest) (*models.UpdateSubscriberRequest, error) {
	if !s.hasAdminPrivilege(request.Email) {
		return nil, errors.New("no admin privilege")
	}
	err := s.dynamodbAdapter.UpdateSubscriberBasic(request.EditedSubscriber)
	if err != nil {
		return nil, err
	}
	return request.EditedSubscriber, nil
}

func (s *Store) AdminSearchSubscribers(subscriberFilter string) ([]*models.Subscriber, error) {
	log.Infof("SearchSubscribers(): %s", subscriberFilter)
	subscribers, err := s.dynamodbAdapter.SearchSubscribersByKey(subscriberFilter)
	if err != nil {
		return nil, err
	}
	return subscribers, nil
}

func (s *Store) adminReadSubscribers(token string, pageSize int) (*models.AdminReadSubscribersResponse, error) {
	var exclusiveStartKey map[string]*dynamodb.AttributeValue
	if token != "" && strings.Contains(token, "#") {
		exclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"projectId": {
				S: aws.String(strings.Split(token, "#")[0]),
			},
			"email": {
				S: aws.String(strings.Split(token, "#")[1]),
			},
		}
	}

	subscribers, nextExclusiveStartKey, err := s.dynamodbAdapter.ScanSubscriber(pageSize, exclusiveStartKey)
	if err != nil {
		return nil, err
	}
	resp := &models.AdminReadSubscribersResponse{
		Subscribers: subscribers,
	}
	if len(nextExclusiveStartKey) > 0 {
		resp.Token = nextExclusiveStartKey["projectId"].String() + "#" + nextExclusiveStartKey["email"].String()
	}
	return resp, nil
}

func (s *Store) AdminDeleteSubscriber(adminEmail, projectId, email string) error {
	log.Infof("AdminDeleteSubscriber(): %s, %s, %s", adminEmail, projectId, email)
	if !s.hasAdminPrivilege(adminEmail) {
		return errors.New("no admin privilege")
	}
	return s.dynamodbAdapter.DeleteSubscriber(projectId, email)
}

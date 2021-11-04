package store

import (
	"email_action/models"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (s *Store) CreateSubscriber(projectId, email, firstName, lastName string) (*models.Subscriber, error) {
	log.Infof("CreateSubscriber(): %s, %s", email, projectId)
	subscriber := models.NewSubscriber(projectId, email, firstName, lastName)
	err := s.dynamodbAdapter.CreateSubscriberIfNotExist(subscriber)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				log.Warnf("CreateSubscriber(): email %s is already subscribed to project %s", email, projectId)
				return subscriber, nil
			}
		}
		return nil, err
	}
	return subscriber, nil
}

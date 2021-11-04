package store

import (
	"email_action/models"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (s *Store) AdminReadUsers(adminEmail, token string, pageSize int) (
	*models.AdminReadUsersResponse, error) {
	log.Infof("AdminReadUsers(): %s, %s, %d", adminEmail, token, pageSize)
	if !s.hasAdminPrivilege(adminEmail) {
		return nil, errors.New("no admin privilege")
	}
	return s.adminReadUserAccountInfos(token, pageSize)
}

func (s *Store) AdminUpdateUser(request *models.AdminUpdateUserRequest) (*models.UserAccountInfo, error) {
	if !s.hasAdminPrivilege(request.Email) {
		return nil, errors.New("no admin privilege")
	}
	userAccountInfo, err := s.dynamodbAdapter.UpdateEditedUserInfo(request.EditedUserInfo)
	if err != nil {
		return nil, err
	}
	return userAccountInfo, nil
}

func (s *Store) adminReadUserAccountInfos(token string, pageSize int) (*models.AdminReadUsersResponse, error) {
	var exclusiveStartKey map[string]*dynamodb.AttributeValue
	if token != "" {
		exclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(token),
			},
		}
	}

	accountInfos, nextExclusiveStartKey, err := s.dynamodbAdapter.ScanUser(pageSize, exclusiveStartKey)
	if err != nil {
		return nil, err
	}
	resp := &models.AdminReadUsersResponse{
		UserAccountInfos: accountInfos,
	}
	if len(nextExclusiveStartKey) > 0 {
		resp.Token = nextExclusiveStartKey["email"].String()
	}
	return resp, nil
}

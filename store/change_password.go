package store

import (
	"email_action/models"
	"errors"
)

func (s *Store) ChangePassword(request models.UpdatePasswordRequest) error {
	log.Infof("ChangePassword(): %s", request.Email)
	authReq := &models.AuthRequest{
		Email:    request.Email,
		Password: request.CurrentPassword,
	}
	_, err := s.dynamodbAdapter.AuthUser(authReq)
	if err != nil {
		return errors.New("incorrect password")
	}
	hashedPassword, err := generateHashPassword(request.NewPassword)
	if err != nil {
		return err
	}
	if request.NewPassword != request.NewPasswordConfirm {
		return errors.New("password does not match")
	}
	err = s.dynamodbAdapter.ChangePassword(request.Email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

package store

import (
	"email_action/models"
	"errors"
)

func (s *Store) UpdatePassword(request models.UpdatePasswordTokenRequest) error {
	log.Infof("UpdatePassword(): %s", request.Token)
	hashedPassword, err := generateHashPassword(request.NewPassword)
	if err != nil {
		return err
	}
	if request.NewPassword != request.NewPasswordConfirm {
		return errors.New("password does not match")
	}

	resetToken, err := s.dynamodbAdapter.ReadPasswordReset(request.Token)
	if err != nil {
		return err
	}
	err = s.dynamodbAdapter.UpdatePassword(resetToken.Email, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

package store

import (
	"errors"
	"time"
)

func (s *Store) ReadPasswordReset(token string) error {
	log.Infof("ReadPasswordReset(): %s", token)
	passwordToken, err := s.dynamodbAdapter.ReadPasswordReset(token)
	if err != nil {
		return err
	}
	if passwordToken.CreatedAt/1000 < time.Now().Unix()-int64(time.Duration(time.Hour).Seconds()) {
		// token expires after 1 hour
		return errors.New("token expired after 1 hour")
	}
	return nil
}

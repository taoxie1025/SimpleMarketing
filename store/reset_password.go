package store

import (
	"github.com/astaxie/beego"
	"time"
)

func (s *Store) ResetPassword(email string) (string, error) {
	log.Infof("ResetPassword(): %s", email)
	token, err := s.dynamodbAdapter.ResetPassword(email, time.Now().UnixNano()/int64(time.Millisecond))
	if err != nil {
		return "", err
	}
	return getResetPasswordLink(token), nil
}

func getResetPasswordLink(token string) string {
	httpAddr := beego.AppConfig.String("httpaddr")
	return "https://" + httpAddr + "/recovery?token=" + token
}

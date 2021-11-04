package controllers

import (
	"email_action/logging"
	"errors"
)

var (
	log             = logging.NewZapLogger()
	unauthorizedErr = errors.New("unauthorized request")
)

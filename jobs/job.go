package jobs

import "email_action/logging"

var (
	log = logging.NewZapLogger()
)

type JobInterface interface {
	Run()
}

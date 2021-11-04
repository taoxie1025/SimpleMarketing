package jobs

import (
	"email_action/mail"
	"email_action/models"
	"email_action/store"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

var (
	UserBatchSize = 50
)

type ClearFreePlanEmailQuotaJob struct {
	JobInterface
	store  *store.Store
	mailer *mail.SesAdapter
}

func NewClearFreePlanEmailQuotaJob(store *store.Store, mailer *mail.SesAdapter) *ClearFreePlanEmailQuotaJob {
	return &ClearFreePlanEmailQuotaJob{
		store:  store,
		mailer: mailer,
	}
}

func (j *ClearFreePlanEmailQuotaJob) Run() {
	log.Infof("running ClearFreePlanEmailQuotaJob...")
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue
	for {
		userAccountInfos, nextKey, err := j.store.GetDynamodbAdapter().ScanUser(UserBatchSize, lastEvaluatedKey)
		if err != nil {
			log.Warnf("ClearFreePlanEmailQuotaJob(): error = ", err)
			continue
		}
		for _, userAccountInfo := range userAccountInfos {
			if !userAccountInfo.IsBlock && userAccountInfo.SubscriptionPlan == models.SubscriptionPlanFree &&
				userAccountInfo.LastClearCycleTime < time.Now().AddDate(0, 0, -30).UnixNano()/int64(time.Millisecond) {
				go j.updateFreePlanUserEmailQuota(userAccountInfo.Email)
			}
		}

		lastEvaluatedKey = nextKey
		if len(lastEvaluatedKey) == 0 {
			return
		}
	}
}

func (j *ClearFreePlanEmailQuotaJob) updateFreePlanUserEmailQuota(email string) error {
	log.Infof("updateFreePlanUserEmailQuota(): %s", email)
	return j.store.GetDynamodbAdapter().ClearEmailQuota(email)
}

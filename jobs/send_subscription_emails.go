package jobs

import (
	"email_action/mail"
	"email_action/models"
	"email_action/store"
	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

var (
	ProjectBatchSize    = 100
	SubscriberBatchSize = 200
	SendEmailInterval   = 1000 // in millisecond
)

type SendSubscriptionEmailsJob struct {
	JobInterface
	store  *store.Store
	mailer *mail.SesAdapter
}

func NewSendSubscriptionEmailsJob(store *store.Store, mailer *mail.SesAdapter) *SendSubscriptionEmailsJob {
	return &SendSubscriptionEmailsJob{
		store:  store,
		mailer: mailer,
	}
}

func (j *SendSubscriptionEmailsJob) Run() {
	log.Infof("running SendSubscriptionEmailsJob...")
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue
	for {
		projects, nextKey, err := j.store.GetDynamodbAdapter().ScanProject(ProjectBatchSize, lastEvaluatedKey)
		if err != nil {
			log.Warnf("SendSubscriptionEmailsJob(): error = ", err)
			continue
		}
		for _, project := range projects {
			if project.ProjectState == models.ProjectStateLive {
				start := time.Now()
				numSuccess, err := j.broadcastProject(project)
				if err != nil {
					continue
				}
				end := time.Now()
				project.TotalBroadcastCount += int64(numSuccess)
				project.LastBroadcastCount = int64(numSuccess)
				project.LastBroadcastTimeMs = time.Now().UnixNano() / int64(time.Millisecond)
				project.LastBroadcastDuration = (end.Sub(start)).Milliseconds()
				if numSuccess > 0 {
					go j.updateProjectStat(project)
				}
			}
		}

		lastEvaluatedKey = nextKey
		if len(lastEvaluatedKey) == 0 {
			return
		}
	}
}

func (j *SendSubscriptionEmailsJob) broadcastProject(project *models.Project) (int, error) {
	log.Infof("broadcastProject(): broadcasting project(%s, %s)", project.Email, project.ProjectId)
	token := ""
	numSuccess := 0

	for {
		projectOwner, err := j.store.GetDynamodbAdapter().ReadUser(project.Email)
		if err != nil {
			return 0, err
		}
		totalQuota := beego.AppConfig.DefaultInt64("free_plan_quota", 1000)
		if projectOwner.SubscriptionPriceId == beego.AppConfig.String("pro_plan_price_id") {
			totalQuota = beego.AppConfig.DefaultInt64("pro_plan_quota", 100000)
		} else if projectOwner.SubscriptionPriceId == beego.AppConfig.String("ultra_plan_price_id") {
			totalQuota = beego.AppConfig.DefaultInt64("ultra_plan_quota", 300000)
		}
		quotaLeft := totalQuota - projectOwner.EmailUsageInCycle

		if !projectOwner.IsBlock && quotaLeft > 0 {
			subscribers, nextToken, err := j.store.GetDynamodbAdapter().ReadSubscribers(project.ProjectId, token, SubscriberBatchSize)
			if err != nil {
				log.Warnf("broadcastProject(): error = ", err)
				return numSuccess, err
			}

			for _, subscriber := range subscribers {
				if quotaLeft > 0 {
					count, err := j.broadcastToSubscriber(subscriber, project)
					if err != nil {
						continue
					}
					numSuccess += count
					time.Sleep(time.Duration(SendEmailInterval))
					quotaLeft--
				} else {
					break
				}
			}

			if numSuccess > 0 {
				j.store.AddEmailUsageInCycle(projectOwner.Email, int64(numSuccess)) // Do not use go routine here because concurrent updates are possible
			}

			if quotaLeft == 0 {
				// TODO: when quotaLeft == 0, send notification to project owner
				return numSuccess, nil
			}

			token = nextToken
			if token == "EOF" {
				return numSuccess, nil
			}
		}
	}
}

func (j *SendSubscriptionEmailsJob) broadcastToSubscriber(subscriber *models.Subscriber, project *models.Project) (int, error) {
	log.Infof("broadcastToSubscriber(): broadcasting to subscriber(%v) for project %s", subscriber, project.ProjectId)

	now := time.Now().UnixNano() / int64(time.Millisecond)
	if subscriber.LastBroadcastTimeMs+project.Interval <= now+time.Hour.Milliseconds() && subscriber.IsEnabled { // + 1 hour to account for processing time
		nextArticlePtr := subscriber.ArticleCursor
		if nextArticlePtr < len(project.ArticleIds) {
			nextArticleId := project.ArticleIds[nextArticlePtr]
			article, err := j.store.GetDynamodbAdapter().ReadArticle(project.Email, nextArticleId)
			if err != nil {
				return 0, err
			}
			if article.IsLive {
				message := article.HtmlBody + j.store.GetUnsubscribeLink(project.ProjectId, project.Name)
				err = j.mailer.SendEmail(project.OutgoingEmail, subscriber.Email, article.Title, message, message)
				if err == nil {
					subscriber.LastBroadcastTimeMs = now
					subscriber.ArticleCursor = nextArticlePtr + 1
					go j.updateSubscriber(subscriber)
					return 1, nil
				} else {
					return 0, err
				}
			} else {
				subscriber.ArticleCursor = nextArticlePtr + 1
				go j.updateSubscriber(subscriber)
				return 0, nil
			}
		}
	}
	return 0, nil
}

func (j *SendSubscriptionEmailsJob) updateProjectStat(project *models.Project) error {
	log.Infof("updateProjectStat(): updating project(%s, %s)", project.Email, project.ProjectId)
	return j.store.GetDynamodbAdapter().UpdateProjectBroadcastStat(project)
}

func (j *SendSubscriptionEmailsJob) updateSubscriber(subscriber *models.Subscriber) error {
	log.Infof("updateSubscriber(): updating subscriber(%s, %s)", subscriber.Email, subscriber.ProjectId)
	return j.store.GetDynamodbAdapter().UpdateSubscriberStat(subscriber)
}

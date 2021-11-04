package store

import (
	"email_action/logging"
	"email_action/models"
	"email_action/store/dynamodb"
	"email_action/store/s3"
	stripe_helper "email_action/store/stripe"
	"github.com/stripe/stripe-go/v72"
)

var (
	log            = logging.NewZapLogger()
	BATCH_GET_SIZE = 50
)

type StoreApi interface {
	AuthUser(authReq *models.AuthRequest) (*models.Claims, string, error)
	CreateArticle(email, projectId, title, htmlBody, textBody string) (*models.Article, error)
	CreateProject(email, name, intro, outgoingEmail, avatarUrl, backgroundImageUrl, author string, interval int64, subscriptionType int) (*models.Project, error)
	CreateCheckoutSession(stripeCustomerId, priceId, plan, email, coupon, successURL, cancelURL string) (*stripe.CheckoutSession, error)
	CreateSubscriber(projectId, email, firstName, lastName string) (*models.Subscriber, error)
	CreateUser(email, password, firstName, lastName, ip string) (*models.User, error)
	DeleteArticle(email, projectId, articleId string) error
	DeleteFile(userEmail, fileName string) error
	DeleteProject(email, projectId string) error
	DeleteSubscriber(projectId, email string) error
	ListFiles(userEmail string) ([]*models.File, error)
	ProcessStripeEvents(reqBody []byte, header string) error
	UpdateSubscriptionPlan(email string, targetPlan int, emailUsageInCycle int64, subscriptionId, subscriptionPriceId string, isBlock bool) error
	UpdateAccountStatus(email string, isBlock bool, paymentStatus int) error
	ReadArticle(email, articleId string) (*models.Article, error)
	ReadArticles(email string, articleIds []string) ([]*models.Article, error)
	ReadProject(email, projectId string) (*models.Project, error)
	ReadProjectBrief(projectId string) (*models.ProjectBrief, error)
	ReadProjects(email string, projectIds []string) ([]*models.Project, error)
	ReadSubscribers(projectId, token string, pageSize int) (*models.ReadSubscribersResponse, error)
	ReadUser(email string) (*models.User, error)
	SaveArticle(article *models.Article) (*models.Article, error)
	SearchSubscribers(projectId, emailFilter string) ([]*models.Subscriber, error)
	UpdateArticle(article *models.Article) (*models.Article, error)
	SaveProject(project *models.Project) (*models.Project, error)
	UpdateSubscriber(req *models.UpdateSubscriberRequest) (*models.Subscriber, error)
	UpdateUserInfo(user *models.User) (*models.User, error)
	UploadFile(file []byte, userEmail, fileName string) (string, error)
	ChangePassword(request models.UpdatePasswordRequest) error
	ResetPassword(email string) (string, error)
	ReadPasswordReset(token string) error
	UpdatePassword(request models.UpdatePasswordTokenRequest) error
	AddEmailUsageInCycle(email string, count int64) error
	GetUnsubscribeLink(projectId, projectName string) string
	ReadCoupon(couponId string) (*stripe.Coupon, error)
	ReadTickets(email string, token string, pageSize int) (*models.ReadTicketsResponse, error)
	CreateTicket(email, title, body string, ticketType int) (*models.Ticket, error)
	CreateComment(email, ticketId string, comment *models.Comment) (*models.Comment, error)
	UpdateTicket(ticket *models.Ticket) (*models.Ticket, error)
	DeleteTicket(email, ticketId string) error
	AdminReadTickets(adminEmail string, token string, pageSize int) (*models.ReadTicketsResponse, error)
	AdminReadUser(adminEmail, token string, pageSize int)
	SearchUsers(emailFilter string) ([]*models.UserAccountInfo, error)
	SearchTickets(ticketFilter string) ([]*models.Ticket, error)
	AdminUpdateUser(request *models.AdminUpdateUserRequest) (*models.UserAccountInfo, error)
}

type Store struct {
	dynamodbAdapter *dynamodb.DynamodbAdapter
	s3Adapter       *s3.S3Adapter
	stripeAdapter   *stripe_helper.StripeAdapter
	jwtSecret       []byte
}

func NewStore(env string, jwtSecret string) *Store {
	log.Infof("NewStore(): env = %s", env)
	dynamodbAdapter := dynamodb.NewDynamodbAdapter(env)
	s3Adapter := s3.NewS3Adapter(env)
	stripeAdapter := stripe_helper.NewStripeAdapter()
	return &Store{
		dynamodbAdapter: dynamodbAdapter,
		jwtSecret:       []byte(jwtSecret),
		s3Adapter:       s3Adapter,
		stripeAdapter:   stripeAdapter,
	}
}

func (s *Store) GetDynamodbAdapter() *dynamodb.DynamodbAdapter {
	return s.dynamodbAdapter
}

func (s *Store) GetS3Adapter() *s3.S3Adapter {
	return s.s3Adapter
}

func (s *Store) GetStripeAdapter() *stripe_helper.StripeAdapter {
	return s.stripeAdapter
}

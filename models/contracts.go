package models

type NewUserRequest struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
}

type UserInfo struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
}

type UpdateUserRequest struct {
	Email    string   `json:"email"`
	UserInfo UserInfo `json:"userInfo"`
}

type UpdatePasswordRequest struct {
	Email              string `json:"email"`
	CurrentPassword    string `json:"currentPassword"`
	NewPassword        string `json:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm"`
}

type UpdatePasswordTokenRequest struct {
	Token              string `json:"token"`
	NewPassword        string `json:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm"`
}

type NewProjectRequest struct {
	Email              string `json:"email"`
	Name               string `json:"name"`
	Intro              string `json:"intro"`
	BackgroundImageUrl string `json:"backgroundImageUrl"`
	AvatarUrl          string `json:"avatarUrl"`
	OutgoingEmail      string `json:"outgoingEmail"`
	Interval           int64  `json:"interval"`
	Author             string `json:"author"`
	SubscriptionType   int    `json:"subscriptionType"`
}

type UpdateProjectRequest struct {
	Email   string  `json:"email"`
	Project Project `json:"project"`
}

type ReadProjectRequest struct {
	Email string `json:"email"`
}

type NewArticleRequest struct {
	Email     string `json:"email"`
	ProjectId string `json:"projectId"`
	IsLive    bool   `json:"isLive"`
	Title     string `json:"title"`
	HtmlBody  string `json:"htmlBody"`
	TextBody  string `json:"textBody"`
}

type UpdateArticleRequest struct {
	Email   string  `json:"email"`
	Article Article `json:"article"`
}

type NewSubscriberRequest struct {
	ProjectId string `json:"projectId"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateSubscriberRequest struct {
	ProjectId     string `json:"projectId"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	IsEnabled     bool   `json:"isEnabled"`
	ArticleCursor int    `json:"articleCursor"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type Claims struct {
	Email     string `json:"email"`
	UserScope int    `json:"userScope"`
	Exp       int64  `json:"exp"`
	Iat       int64  `json:"iat"`
}

func (c *Claims) Valid() error {
	return nil
}

type AuthResponse struct {
	Token  string  `json:"token"`
	Claims *Claims `json:"claims"`
}

type BatchReadProjectRequest struct {
	Email      string   `json:"email"`
	ProjectIds []string `json:"projectIds"`
}

type BatchReadArticleRequest struct {
	Email      string   `json:"email"`
	ArticleIds []string `json:"articleIds"`
}

type File struct {
	Url          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type EmailVerificationResult struct {
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
}

type ContactUsRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
}

type SendEmailRequest struct {
	Email       string `json:"email"` // userId
	From        string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
}

type ReadSubscribersRequest struct {
	ProjectId string `json:"projectId"`
	Email     string `json:"email"`
	PageSize  int    `json:"pageSize"`
	Token     string `json:"token"`
}

type ReadSubscribersResponse struct {
	Token       string        `json:"token"`
	Subscribers []*Subscriber `json:"subscribers"`
}

type NewCheckoutSession struct {
	PriceId          string `json:"priceId"`
	Plan             string `json:"plan"`
	Email            string `json:"email"`
	SuccessUrl       string `json:"successUrl"`
	CancelUrl        string `json:"cancelUrl"`
	StripeCustomerId string `json:"stripeCustomerId"`
	Coupon           string `json:"coupon"`
}

type UpdateSubscriptionRequest struct {
	PriceId string `json:"priceId"`
}

type NewTicketRequest struct {
	Email       string `json:"email"`
	TicketType  int    `json:"ticketType"`
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Body        string `json:"body"`
}

type UpdateTicketRequest struct {
	Email  string  `json:"email"`
	Ticket *Ticket `json:"ticket"`
}

type NewCommentRequest struct {
	Email    string `json:"email"`
	TicketId string `json:"ticketId"`
	Name     string `json:"name"`
	Body     string `json:"body"`
}

type ReadTicketsRequest struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	PageSize int    `json:"pageSize"`
}

type ReadTicketsResponse struct {
	Tickets []*Ticket `json:"tickets"`
	Token   string    `json:"token"`
}

type AdminReadTicketsRequest struct {
	Email    string `json:"email"` // admin email
	Token    string `json:"token"`
	PageSize int    `json:"pageSize"`
}

type AdminReadTicketsResponse struct {
	Token   string    `json:"token"`
	Tickets []*Ticket `json:"tickets"`
}

type AdminReadUsersRequest struct {
	Email    string `json:"email"` // admin email
	Token    string `json:"token"`
	PageSize int    `json:"pageSize"`
}

type AdminReadUsersResponse struct {
	Token            string             `json:"token"`
	UserAccountInfos []*UserAccountInfo `json:"userAccountInfos"`
}

type EditedUserInfo struct {
	Email             string `json:"email"`
	IsBlock           bool   `json:"isBlock"`
	UserScope         int    `json:"userScope"`
	SubscriptionPlan  int    `json:"subscriptionPlan"`
	EmailUsageInCycle int64  `json:"emailUsageInCycle"`
	StripeCustomerId  string `json:"stripeCustomerId"`
	PaymentStatus     int    `json:"paymentStatus"`
}

type AdminUpdateUserRequest struct {
	Email          string          `json:"email"` // admin email
	EditedUserInfo *EditedUserInfo `json:"editedUserInfo"`
}

type AdminReadSubscribersRequest struct {
	Email    string `json:"email"` // admin email
	Token    string `json:"token"`
	PageSize int    `json:"pageSize"`
}

type AdminReadSubscribersResponse struct {
	Token       string        `json:"token"`
	Subscribers []*Subscriber `json:"subscribers"`
}

type AdminUpdateSubscriberRequest struct {
	Email            string                   `json:"email"` // admin email
	EditedSubscriber *UpdateSubscriberRequest `json:"editedSubscriber"`
}

type AdminDeleteSubscriberRequest struct {
	Email           string `json:"email"` // admin email
	ProjectId       string `json:"projectId"`
	SubscriberEmail string `json:"subscriberEmail"`
}

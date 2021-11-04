package models

type UserAccountInfo struct {
	Email               string `json:"email"`
	StripeCustomerId    string `json:"stripeCustomerId"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	EmailUsageInCycle   int64  `json:"emailUsageInCycle"`
	SubscriptionPlan    int    `json:"subscriptionPlan"`
	UserScope           int    `json:"userScope"`
	IsPaymentError      bool   `json:"isPaymentError"`
	IsBlock             bool   `json:"isBlock"`
	PaymentStatus       int    `json:"paymentStatus"`
	SubscriptionId      string `json:"subscriptionId"`
	SubscriptionPriceId string `json:"subscriptionPriceId"`
	Address             string `json:"address"`
	PhoneNumber         string `json:"phoneNumber"`
	LastClearCycleTime  int64  `json:"lastClearCycleTime"`
}

func NewUserAccountInfo(user *User) *UserAccountInfo {
	return &UserAccountInfo{
		Email:               user.Email,
		StripeCustomerId:    user.StripeCustomerId,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		EmailUsageInCycle:   user.EmailUsageInCycle,
		SubscriptionPlan:    user.SubscriptionPlan,
		UserScope:           user.UserScope,
		PaymentStatus:       user.PaymentStatus,
		IsBlock:             user.IsBlock,
		SubscriptionId:      user.SubscriptionId,
		SubscriptionPriceId: user.SubscriptionPriceId,
		Address:             user.Address,
		PhoneNumber:         user.PhoneNumber,
		LastClearCycleTime:  user.LastClearCycleTime,
	}
}

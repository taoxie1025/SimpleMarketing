package models

import (
	"time"
)

const (
	ScopeConsumer = iota
	ScopeAdmin
	ScopeSuperAdmin
)

const (
	SubscriptionPlanFree = iota
	SubscriptionPlanPro
	SubscriptionPlanUltra
)

const (
	PaymentStatusCompleted = iota
	PaymentStatusInProgress
	PaymentStatusFailed
)

type User struct {
	Email               string   `json:"email"` // Email is the primary key
	CreatedAt           int64    `json:"createdAt"`
	ProjectIds          []string `json:"projectIds"`
	LastSeenIp          string   `json:"lastSeenIp"`
	LastSignInTime      int64    `json:"lastSignInTime"`
	FirstName           string   `json:"firstName"`
	LastName            string   `json:"lastName"`
	PasswordHash        string   `json:"passwordHash, omitempty"`
	IsBlock             bool     `json:"isBlock"`
	UserScope           int      `json:"userScope"`
	SubscriptionPlan    int      `json:"subscriptionPlan"`
	EmailUsageInCycle   int64    `json:"emailUsageInCycle"`
	StripeCustomerId    string   `json:"stripeCustomerId"`
	PaymentStatus       int      `json:"paymentStatus"`
	SubscriptionId      string   `json:"subscriptionId"`
	SubscriptionPriceId string   `json:"subscriptionPriceId"`
	Address             string   `json:"address"`
	PhoneNumber         string   `json:"phoneNumber"`
	LastClearCycleTime  int64    `json:"lastClearCycleTime"`
}

func NewUser(email, passwordHash, firstName, lastName, lastSeenIp, stripeId, address, phoneNum string) *User {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	return &User{
		Email:              email,
		PasswordHash:       passwordHash,
		CreatedAt:          currentTime,
		ProjectIds:         []string{},
		LastSeenIp:         lastSeenIp,
		FirstName:          firstName,
		LastName:           lastName,
		IsBlock:            false,
		SubscriptionPlan:   SubscriptionPlanFree,
		EmailUsageInCycle:  0,
		StripeCustomerId:   stripeId,
		Address:            address,
		PhoneNumber:        phoneNum,
		LastClearCycleTime: currentTime,
	}
}

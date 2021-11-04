package models

import (
	"time"
)

const (
	ProjectStateNone = iota
	ProjectStateCreated
	ProjectStatePending
	ProjectStateLive
	ProjectStateDeleted
	ProjectStateSuspended
)

const (
	DefaultSubscription = iota
	RollingSubscription
)

type Project struct {
	Email                 string   `json:"email"`     // Email is the primary key
	ProjectId             string   `json:"projectId"` // ProjectId is the sort key
	ProjectState          int      `json:"projectState"`
	Name                  string   `json:"name"`
	CreatedAt             int64    `json:"createdAt"`
	UpdatedAt             int64    `json:"updatedAt"`
	ArticleIds            []string `json:"articleIds"`
	Interval              int64    `json:"interval"`
	LastBroadcastTimeMs   int64    `json:"lastBroadcastTimeMs"`
	IsBroadcasting        bool     `json:"isBroadcasting"`
	SubscriberCount       int64    `json:"subscriberCount"`
	LastBroadcastCount    int64    `json:"lastBroadcastCount"`
	LastBroadcastDuration int64    `json:"lastBroadcastDuration"`
	TotalBroadcastCount   int64    `json:"totalBroadcastCount"`
	Intro                 string   `json:"intro"`
	BackgroundImageUrl    string   `json:"backgroundImageUrl"`
	AvatarUrl             string   `json:"avatarUrl"`
	OutgoingEmail         string   `json:"outgoingEmail"`
	Author                string   `json:"author"`
	SubscriptionType      int      `json:"subscriptionType"`
}

func NewProject(email, name, intro, outgoingEmail, avatarUrl, backgroundImageUrl, author string, interval int64, subscriptionType int) *Project {
	st := subscriptionType
	if st < DefaultSubscription || st > RollingSubscription {
		st = DefaultSubscription
	}
	return &Project{
		Email:              email,
		ProjectId:          GenerateProjectUUID(),
		OutgoingEmail:      outgoingEmail,
		Name:               name,
		Intro:              intro,
		BackgroundImageUrl: backgroundImageUrl,
		AvatarUrl:          avatarUrl,
		CreatedAt:          time.Now().UnixNano() / int64(time.Millisecond),
		ArticleIds:         []string{},
		Interval:           interval,
		ProjectState:       ProjectStateCreated,
		Author:             author,
		SubscriptionType:   st,
	}
}

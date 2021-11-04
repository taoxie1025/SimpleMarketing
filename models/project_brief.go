package models

type ProjectBrief struct {
	ProjectId           string `json:"projectId"`
	Name                string `json:"name"`
	CreatedAt           int64  `json:"createdAt"`
	Interval            int64  `json:"interval"`
	LastBroadcastTimeMs int64  `json:"lastBroadcastTimeMs"`
	Intro               string `json:"intro"`
	BackgroundImageUrl  string `json:"backgroundImageUrl"`
	AvatarUrl           string `json:"avatarUrl"`
	OutgoingEmail       string `json:"outgoingEmail"`
	Author              string `json:"author"`
	TotalArticles       int    `json:"totalArticles"`
	SubscriptionType    int    `json:"subscriptionType"`
}

package models

import "time"

type Subscriber struct {
	ProjectId           string `json:"projectId"` // ProjectId is the primary key
	Email               string `json:"email"`     // Email is the sort key
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CreatedAt           int64  `json:"createdAt"`
	ArticleCursor       int    `json:"articleCursor"`
	LastBroadcastTimeMs int64  `json:"lastBroadcastTimeMs"`
	IsEnabled           bool   `json:"isEnabled"`
}

func NewSubscriber(projectId, email, firstName, lastName string) *Subscriber {
	return &Subscriber{
		ProjectId:     projectId,
		Email:         email,
		FirstName:     firstName,
		LastName:      lastName,
		IsEnabled:     true,
		ArticleCursor: 0,
		CreatedAt:     time.Now().UnixNano() / int64(time.Millisecond),
	}
}

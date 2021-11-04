package models

type ResetPasswordToken struct {
	Token     string `json:"token"` // Token is the primary key
	Email     string `json:"email"`
	CreatedAt int64  `json:"createdAt"` // in millisecond
	ExpiredAt int64  `json:"expiredAt"` // in second, dynamodb TTL attribute
}

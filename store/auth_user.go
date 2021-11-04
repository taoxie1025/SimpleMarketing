package store

import (
	"email_action/models"
	"github.com/dgrijalva/jwt-go"
)

func (s *Store) AuthUser(authReq *models.AuthRequest) (*models.Claims, string, error) {
	log.Infof("AuthUser():")
	claims, err := s.dynamodbAdapter.AuthUser(authReq)
	if err != nil {
		return nil, "", err
	}
	token, err := s.generateJwt(claims)
	if err != nil {
		return nil, "", err
	}
	return claims, token, nil
}

func (s *Store) generateJwt(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		log.Infof("generateJwt(): failed to sign, error = %v", err)
	}
	return signedString, err
}

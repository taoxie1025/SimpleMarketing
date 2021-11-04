package store

import (
	"email_action/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *Store) CreateUser(email, password, firstName, lastName, ip, address, phoneNum string) (*models.User, error) {
	log.Infof("CreateUser(): %s, %s, %s, %s", email, firstName, lastName, ip)
	hashedPassword, err := generateHashPassword(password)
	if err != nil {
		return nil, err
	}
	customer, err := s.GetStripeAdapter().CreateCustomer(email, firstName, lastName)
	if err != nil {
		return nil, err
	}
	user := models.NewUser(email, hashedPassword, firstName, lastName, ip, customer.ID, address, phoneNum)
	err = s.dynamodbAdapter.CreateUserIfNotExist(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func generateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("generateHashPassword(): error = %v", err)
		return "", err
	}
	return string(bytes), err
}

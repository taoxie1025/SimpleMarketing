package store

import (
	"email_action/models"
	"errors"
)

func (s *Store) AdminReadTickets(adminEmail, token string, pageSize int) (
	*models.AdminReadTicketsResponse, error) {
	log.Infof("AdminReadTickets(): %s, %s, %d", adminEmail, token, pageSize)
	if !s.hasAdminPrivilege(adminEmail) {
		return nil, errors.New("no admin privilege")
	}
	return s.adminReadTickets(token, pageSize)
}

func (s *Store) hasAdminPrivilege(email string) bool {
	user, err := s.dynamodbAdapter.ReadUser(email)
	if err != nil {
		return false
	}
	if user.UserScope == models.ScopeAdmin || user.UserScope == models.ScopeSuperAdmin {
		return true
	}
	return false
}

func (s *Store) adminReadTickets(token string, pageSize int) (*models.AdminReadTicketsResponse, error) {
	tickets, nextToken, err := s.dynamodbAdapter.ScanTickets(pageSize, token)
	if err != nil {
		return nil, err
	}
	resp := &models.AdminReadTicketsResponse{
		Token:   nextToken,
		Tickets: tickets,
	}
	return resp, nil
}

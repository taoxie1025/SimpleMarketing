package store

import (
	"email_action/models"
)

func (s *Store) ReadTickets(email string, token string, pageSize int) (*models.ReadTicketsResponse, error) {
	log.Infof("ReadTickets(): %s, %s, %d", email, token, pageSize)
	return s.userReadTickets(email, token, pageSize)
}

func (s *Store) userReadTickets(email string, token string, pageSize int) (*models.ReadTicketsResponse, error) {
	tickets, nextToken, err := s.dynamodbAdapter.ReadTickets(email, token, pageSize)
	if err != nil {
		return nil, err
	}
	resp := &models.ReadTicketsResponse{
		Token:   nextToken,
		Tickets: tickets,
	}
	return resp, nil
}

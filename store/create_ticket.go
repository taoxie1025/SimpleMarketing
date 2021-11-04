package store

import (
	"email_action/models"
)

func (s *Store) CreateTicket(email, projectId, projectName, name, title, body string, ticketType int) (*models.Ticket, error) {
	log.Infof("CreateTicket(): %s, %s, %d", email, projectId, ticketType)
	ticket := models.NewTicket(email, projectId, projectName, name, title, body, ticketType)
	err := s.dynamodbAdapter.CreateTicket(ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

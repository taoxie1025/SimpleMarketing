package store

import (
	"email_action/models"
)

func (s *Store) UpdateTicket(ticket *models.Ticket) (*models.Ticket, error) {
	log.Infof("UpdateTicket(): %s, %d", ticket.TicketId, ticket.Email)
	err := s.dynamodbAdapter.UpdateTicket(ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

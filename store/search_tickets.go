package store

import "email_action/models"

func (s *Store) SearchTickets(ticketFilter string) ([]*models.Ticket, error) {
	log.Infof("SearchTickets(): %s", ticketFilter)
	tickets, err := s.dynamodbAdapter.SearchTickets(ticketFilter)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

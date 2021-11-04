package store

func (s *Store) DeleteTicket(email, ticketId string) error {
	log.Infof("DeleteTicket(): %s, %d", email, ticketId)
	return s.dynamodbAdapter.DeleteTicket(email, ticketId)
}

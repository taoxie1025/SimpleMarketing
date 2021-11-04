package store

func (s *Store) DeleteSubscriber(projectId, email string) error {
	log.Infof("DeleteSubscriber(): %s, %s", projectId, email)
	return s.dynamodbAdapter.DeleteSubscriber(projectId, email)
}

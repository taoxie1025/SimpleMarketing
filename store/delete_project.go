package store

func (s *Store) DeleteProject(email, projectId string) error {
	log.Infof("DeleteProject(): email = %s, projectId = %s", projectId, email)
	return s.dynamodbAdapter.DeleteProject(email, projectId)
}

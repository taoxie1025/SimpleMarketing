package store

func (s *Store) DeleteArticle(email, projectId, articleId string) error {
	log.Infof("DeleteArticle(): %s, %s, %s", email, projectId, articleId)
	return s.dynamodbAdapter.DeleteArticle(email, projectId, articleId)
}

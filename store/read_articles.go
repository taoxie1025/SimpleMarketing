package store

import (
	"email_action/models"
)

func (s *Store) ReadArticles(email string, articleIds []string) ([]*models.Article, error) {
	log.Infof("ReadArticles(): %s, %v", email, articleIds)
	var idPairs [][]string
	articles := []*models.Article{}
	for _, projectId := range articleIds {
		id := []string{email, projectId}
		idPairs = append(idPairs, id)
	}

	var idPairChunks [][][]string
	for i := 0; i < len(idPairs); i += BATCH_GET_SIZE {
		end := i + BATCH_GET_SIZE
		if end > len(idPairs) {
			end = len(idPairs)
		}
		idPairChunks = append(idPairChunks, idPairs[i:end])
	}
	for _, idPairChunk := range idPairChunks {
		moreArticles, err := s.dynamodbAdapter.BatchReadArticles(idPairChunk)
		if err != nil {
			return nil, err
		}
		articles = append(articles, moreArticles...)
	}
	sorted := sortArticles(articleIds, articles)
	return sorted, nil
}

// sortArticles sorts articles based on original articleIds order
func sortArticles(articleIds []string, original []*models.Article) []*models.Article {
	idMap := make(map[string]*models.Article)
	for _, article := range original {
		idMap[article.ArticleId] = article
	}
	sorted := []*models.Article{}
	for _, articleId := range articleIds {
		if _, ok := idMap[articleId]; ok {
			sorted = append(sorted, idMap[articleId])
		}
	}
	return sorted
}

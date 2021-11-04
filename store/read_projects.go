package store

import (
	"email_action/models"
)

func (s *Store) ReadProjects(email string, projectIds []string) ([]*models.Project, error) {
	log.Infof("ReadProjects(): %s, %v", email, projectIds)
	var idPairs [][]string
	projects := []*models.Project{}
	for _, projectId := range projectIds {
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
		moreProjects, err := s.dynamodbAdapter.BatchReadProjects(idPairChunk)
		if err != nil {
			return nil, err
		}
		projects = append(projects, moreProjects...)
	}
	return projects, nil
}

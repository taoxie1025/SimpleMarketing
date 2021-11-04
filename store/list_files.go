package store

import "email_action/models"

func (s *Store) ListFiles(userEmail string) ([]*models.File, error) {
	log.Infof("ListFiles(): userEmail = %s", userEmail)
	urls, err := s.GetS3Adapter().ListFiles(userEmail)
	if err != nil {
		return []*models.File{}, err
	}

	files := []*models.File{}
	for _, url := range urls {
		file := &models.File{
			Url:          url,
			ThumbnailUrl: url + "-thumbnail",
		}
		files = append(files, file)
	}
	return files, nil
}

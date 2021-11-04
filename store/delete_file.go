package store

func (s *Store) DeleteFile(userEmail, fileName string) error {
	log.Infof("DeleteFile(): userEmail = %s, fileName = %s", userEmail, fileName)
	key := userEmail + "/" + fileName
	return s.GetS3Adapter().DeleteFile(key)
}

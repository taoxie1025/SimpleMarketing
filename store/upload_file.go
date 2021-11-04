package store

func (s *Store) UploadFile(file []byte, userEmail, fileName string) (string, error) {
	log.Infof("UploadFile(): size = %d, userEmail = %s, fileName = %s", len(file), userEmail, fileName)
	key := userEmail + "/" + fileName
	return s.GetS3Adapter().UploadFile(file, key)
}

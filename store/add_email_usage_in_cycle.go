package store

func (s *Store) AddEmailUsageInCycle(email string, count int64) error {
	log.Infof("AddEmailUsageInCycle(): adding emailUsageInCycle(%s, %d)", email, count)
	return s.GetDynamodbAdapter().AddEmailUsageInCycle(email, count)
}

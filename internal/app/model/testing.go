package model
 
type testStorage struct {
	orders map[string]string
}

func (s *testStorage) Open(dbUrl string) error {
	return nil
}

func (s *testStorage) Load(key string) (string, bool) {
	val, ok := s.orders[key]
	return val, ok
}

func (s *testStorage) Store(key string, value string) error {
	s.orders[key] = value
	return nil
}

func (s *testStorage) Close() {

}

func NewTestStorage() *testStorage {
	return &testStorage{
		orders: make(map[string]string),
	}
}
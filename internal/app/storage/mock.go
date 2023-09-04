package storage

type MockStorage map[string]string

func (m MockStorage) Save(id string, value string) {
	m[id] = value
}

func (m MockStorage) Find(id string) (string, bool) {
	val, ok := m[id]
	return val, ok
}

var MockStorageImpl MockStorage = make(map[string]string)

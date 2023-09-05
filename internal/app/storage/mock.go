package storage

type mockStorage map[string]string

func (m mockStorage) Save(id string, value string) {
	m[id] = value
}

func (m mockStorage) Find(id string) (string, bool) {
	val, ok := m[id]
	return val, ok
}

var mockStorageImpl mockStorage = make(map[string]string)

func GetStorage() Repository {
	return mockStorageImpl
}

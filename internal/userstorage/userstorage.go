package userstorage

// Repository Репозиторий для работы с хранилищем пользователей.
type Repository interface {

	// Save Сохранить пользователя.
	Save(string) error

	// Find Найти пользователя.
	Find(string) (string, bool, error)
}

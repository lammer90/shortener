package userstorage

type Repository interface {
	Save(string) error
	Find(string) (string, bool, error)
}

package storage

type Repository interface {
	Save(string, string) error
	Find(string) (string, bool, error)
}

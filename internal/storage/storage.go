package storage

type Repository interface {
	Save(string, string)
	Find(string) (string, bool)
}

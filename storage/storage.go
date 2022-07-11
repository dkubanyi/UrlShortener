package storage

type Service interface {
	Save(string) (string, error)
	Load(string) (string, error)
	Close() error
}

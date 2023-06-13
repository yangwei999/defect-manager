package backend

type CveBackend interface {
	MaxBulletinID() (int, error)
	PublishedDefects() ([]string, error)
}

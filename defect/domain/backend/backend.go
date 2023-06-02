package backend

type CveBackend interface {
	MaxBulletinID() (int, error)
	IsDefectPublished([]string) (map[string]bool, error)
}

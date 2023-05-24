package repositoryimpl

type dbimpl interface {
	FirstOrCreate(filter, result interface{}) error
	GetRecord(filter, result interface{}) error
	UpdateRecord(filter, update interface{}) error
}

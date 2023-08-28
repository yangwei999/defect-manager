package repositoryimpl

import (
	postgres "github.com/opensourceways/server-common-lib/postgre"
)

type dbimpl interface {
	GetRecord(filter, result interface{}) error
	Insert(result interface{}) error
	UpdateRecord(filter, update interface{}) error
	GetRecords(
		filter []postgres.ColumnFilter, result interface{}, p postgres.Pagination, sort []postgres.SortByColumn,
	) error

	AutoMigrate(dst interface{}) error

	IsRowNotFound(error) bool
	IsRowExists(error) bool
}

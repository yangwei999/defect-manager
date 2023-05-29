package repositoryimpl

import (
	"github.com/opensourceways/defect-manager/common/infrastructure/postgres"
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

package postgres

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	errRowExists   = errors.New("row exists")
	errRowNotFound = errors.New("row not found")
)

type SortByColumn struct {
	Column string
	Ascend bool
}

func (s SortByColumn) order() string {
	v := " ASC"
	if !s.Ascend {
		v = " DESC"
	}
	return s.Column + v
}

type Pagination struct {
	PageNum      int
	CountPerPage int
}

func (p Pagination) pagination() (limit, offset int) {
	limit = p.CountPerPage

	if limit > 0 && p.PageNum > 0 {
		offset = (p.PageNum - 1) * limit
	}

	return
}

type ColumnFilter struct {
	column string
	symbol string
	value  interface{}
}

func (q *ColumnFilter) condition() string {
	return fmt.Sprintf("%s %s ?", q.column, q.symbol)
}

func NewEqualFilter(column string, value interface{}) ColumnFilter {
	return ColumnFilter{
		column: column,
		symbol: "=",
		value:  value,
	}
}

func NewLikeFilter(column string, value string) ColumnFilter {
	return ColumnFilter{
		column: column,
		symbol: "like",
		value:  "%" + value + "%",
	}
}

type DbTable struct {
	name string
}

func NewDBTable(name string) DbTable {
	return DbTable{name: name}
}

func (t DbTable) FirstOrCreate(filter, result interface{}) error {
	query := db.Table(t.name).Where(filter).FirstOrCreate(result)

	if err := query.Error; err != nil {
		return err
	}

	if query.RowsAffected == 0 {
		return errRowExists
	}

	return nil
}

func (t DbTable) Insert(result interface{}) error {
	return db.Table(t.name).Create(result).Error
}

func (t DbTable) FirstOrCreateWithNot(filter, notFilter, result interface{}) error {
	query := db.Table(t.name).
		Where(filter).
		Not(notFilter).
		FirstOrCreate(result)

	if err := query.Error; err != nil {
		return err
	}

	if query.RowsAffected == 0 {
		return errRowExists
	}

	return nil
}

func (t DbTable) GetRecords(
	filter []ColumnFilter, result interface{}, p Pagination, sort []SortByColumn,
) (err error) {
	query := db.Table(t.name)
	for i := range filter {
		query.Where(filter[i].condition(), filter[i].value)
	}

	var orders []string
	for _, v := range sort {
		orders = append(orders, v.order())
	}

	if len(orders) >= 0 {
		query.Order(strings.Join(orders, ","))
	}

	if limit, offset := p.pagination(); limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	err = query.Find(result).Error

	return
}

func (t DbTable) Count(filter []ColumnFilter) (int, error) {
	var total int64
	query := db.Table(t.name)
	for i := range filter {
		query.Where(filter[i].condition(), filter[i].value)
	}

	err := query.Count(&total).Error

	return int(total), err
}

func (t DbTable) GetRecord(filter, result interface{}) error {
	err := db.Table(t.name).Where(filter).First(result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errRowNotFound
	}

	return err
}

func (t DbTable) UpdateRecord(filter, update interface{}) (err error) {
	query := db.Table(t.name).Where(filter).Updates(update)
	if err = query.Error; err != nil {
		return
	}

	if query.RowsAffected == 0 {
		err = errRowNotFound
	}

	return
}

func (t DbTable) ExecSQL(sql string, result interface{}, args ...interface{}) error {
	return db.Exec(sql, args...).Find(result).Error
}

// CreateOrUpdate updates must use primary key [uuid], other fields are invalid.
func (t DbTable) CreateOrUpdate(result interface{}, updates ...string) error {
	return db.Table(t.name).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "uuid"}},
			DoUpdates: clause.AssignmentColumns(updates),
		},
	).Create(result).Error
}

func (t DbTable) DB() *gorm.DB {
	return db
}

func (t DbTable) IsRowNotFound(err error) bool {
	return errors.Is(err, errRowNotFound)
}

func (t DbTable) IsRowExists(err error) bool {
	return errors.Is(err, errRowExists)
}

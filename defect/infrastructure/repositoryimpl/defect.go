package repositoryimpl

import (
	"time"

	"github.com/opensourceways/defect-manager/common/infrastructure/postgres"
	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/repository"
)

type defectImpl struct {
	db dbimpl
}

func NewDefect(cfg *Config) repository.Defect {
	return defectImpl{postgres.NewDBTable(cfg.Table.Defect)}
}

func (d defectImpl) Add(defect domain.Defect) error {
	return nil
}

func (d defectImpl) Save(defect domain.Defect) error {
	return nil
}

func (d defectImpl) FindDefects(time time.Time) ([]domain.Defect, error) {
	return nil, nil
}

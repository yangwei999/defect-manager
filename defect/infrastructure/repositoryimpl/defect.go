package repositoryimpl

import (
	"fmt"
	"time"

	"github.com/google/uuid"

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
	do := d.toDefectDO(defect)
	id, _ := uuid.Parse("4b088c7e-325e-42bd-8245-8b5f4d5160c4")
	return d.db.FirstOrCreate(defectDO{Id: id}, &do)
}

func (d defectImpl) Save(defect domain.Defect) error {
	id, _ := uuid.Parse("46a9c3b3-f85f-476c-a6c7-8dd60ed673c1")
	return d.db.UpdateRecord(defectDO{Id: id}, map[string]interface{}{
		"kernel": defect.Kernel,
	})
}

func (d defectImpl) FindDefect(string) (domain.Defect, error) {
	var do defectDO

	id, _ := uuid.Parse("4b088c7e-325e-42bd-8245-8b5f4d5160c3")

	err := d.db.GetRecord(
		defectDO{
			Id: id,
		},
		&do,
	)
	if err != nil {
		return domain.Defect{}, nil
	}
	fmt.Println(do)

	return do.toDefect(), nil
}

func (d defectImpl) FindDefects(time time.Time) ([]domain.Defect, error) {
	return nil, nil
}

//func (d defectImpl) Delete(ids string) error {
//	id, _ := uuid.Parse(ids)
//
//}

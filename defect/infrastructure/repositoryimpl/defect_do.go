package repositoryimpl

import (
	"github.com/google/uuid"

	"github.com/opensourceways/defect-manager/defect/domain"
)

type defectDO struct {
	Id     uuid.UUID `gorm:"column:uuid;type:uuid" json:"-"`
	Kernel string    `gorm:"column:kernel"         json:"kernel"`
}

func (d defectImpl) toDefectDO(defect domain.Defect) defectDO {
	return defectDO{
		Id:     uuid.New(),
		Kernel: defect.Kernel,
	}
}

func (d defectDO) toDefect(do defectDO) domain.Defect {
	return domain.Defect{
		Kernel: do.Kernel,
	}
}

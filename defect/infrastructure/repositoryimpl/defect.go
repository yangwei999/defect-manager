package repositoryimpl

import (
	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/repository"
)

var instance repository.Defect

func Init() error {
	return nil
}

type defectImpl struct {
}

func Instance() repository.Defect {
	return instance
}

func (impl defectImpl) Save(defect domain.Defect) error {

	return nil
}
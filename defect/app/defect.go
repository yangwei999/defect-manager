package app

import (
	"github.com/opensourceways/defect-manager/defect/domain/repository"
)

type DefectService interface {
	SaveDefect(defect CmdToHandleDefect) error
}

func NewDefectService(repo repository.Defect) *defectService {
	return &defectService{
		repository: repo,
	}
}

type defectService struct {
	repository repository.Defect
}

func (d defectService) SaveDefect(cmd CmdToHandleDefect) error {
	return d.repository.Save(cmd)
}

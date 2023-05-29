package repository

import (
	"time"

	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type OptToFindDefects struct {
	BeginTime time.Time
	Org       string
	Number    []string
	Status    dp.IssueStatus
}

type DefectRepository interface {
	HasDefect(*domain.Defect) (bool, error)
	AddDefect(*domain.Defect) error
	SaveDefect(*domain.Defect) error
	FindDefects(OptToFindDefects) ([]*domain.Defect, error)
}

package repository

import "github.com/opensourceways/defect-manager/defect/domain"

type DefectRepository interface {
	AddDefect(domain.Defect) error
	SaveDefect(domain.Defect) error

	FindDefects(issueNumber []string) ([]domain.Defect, error)
	FindDefectsByTime(time string) ([]domain.Defect, error)
}

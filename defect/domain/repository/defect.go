package repository

import "github.com/opensourceways/defect-manager/defect/domain"

type Defect interface {
	Save(defect domain.Defect) error
}

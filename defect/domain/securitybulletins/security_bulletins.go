package securitybulletins

import "github.com/opensourceways/defect-manager/defect/domain"

type SecurityBulletins interface {
	Generate([]domain.Defect) ([]string, error)
}

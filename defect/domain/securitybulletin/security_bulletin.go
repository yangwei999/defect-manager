package securitybulletin

import "github.com/opensourceways/defect-manager/defect/domain"

type SecurityBulletin interface {
	Generate(domain.SecurityBulletin) (string, error)
}

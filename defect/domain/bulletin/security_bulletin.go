package bulletin

import "github.com/opensourceways/defect-manager/defect/domain"

type Bulletin interface {
	Generate(domain.SecurityBulletin) (string, error)
}

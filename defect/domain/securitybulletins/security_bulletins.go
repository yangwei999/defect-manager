package securitybulletins

import "github.com/opensourceways/defect-manager/defect/domain"

type SecurityBulletins interface {
	ProductTree(component string) domain.ProductTree
	GenerateXML(domain.SecurityBulletins) string
}

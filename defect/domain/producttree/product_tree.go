package producttree

import (
	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type ProductTree interface {
	InitCache()
	CleanCache()
	GetTree(component string, version []dp.SystemVersion) (domain.ProductTree, error)
}

package producttree

import (
	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type ProductTree interface {
	GetTree(component string, version []dp.SystemVersion) (domain.ProductTree, error)
	CleanCache()
}

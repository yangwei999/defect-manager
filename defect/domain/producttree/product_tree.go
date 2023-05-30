package producttree

import "github.com/opensourceways/defect-manager/defect/domain"

type ProductTree interface {
	GetTree(string) (domain.ProductTree, error)
}

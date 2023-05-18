package dp

import "errors"

const (
	productName = "Product Name"
	packageArch = "Package Arch"
)

var validBranchType = map[string]bool{
	productName: true,
	packageArch: true,
}

type branchType string

type BranchType interface {
	BranchType() string
}

func NewBranchType(s string) (BranchType, error) {
	if !validBranchType[s] {
		return nil, errors.New("invalid branch type")
	}

	return branchType(s), nil
}

func (b branchType) BranchType() string {
	return string(b)
}

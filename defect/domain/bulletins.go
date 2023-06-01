package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type SecurityBulletin struct {
	AffectedVersion []dp.SystemVersion
	Identification  string
	Date            string
	Component       string
	ProductTree     ProductTree
	Defects         Defects
}

type ProductTree = map[string][]Product

type Product struct {
	ID       string
	CPE      string
	FullName string
}

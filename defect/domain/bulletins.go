package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type SecurityBulletin struct {
	AffectedVersion []dp.SystemVersion
	Identification  string
	Date            string
	Component       Component
	Defects         Defects
}

type Component struct {
	Name        string
	ProductTree ProductTree
}

type ProductTree struct {
	Branches []Branch
}

type Branch struct {
	Type    string
	Name    dp.Arch
	Product []Product
}

type Product struct {
	ID       string
	CPE      string
	FullName string
}

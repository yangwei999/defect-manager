package domain

import (
	"time"

	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type SecurityBulletin struct {
	AffectedVersion []dp.SystemVersion
	Identification  string
	Date            time.Time
	Component       Component
	Influences      []string
	Vulnerability   []Vulnerability
}

type Component struct {
	Name        string
	ProductTree ProductTree
}

type Vulnerability struct {
	Description     string
	SeverityLevel   dp.SeverityLevel
	AffectedVersion []dp.SystemVersion
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

package domain

import (
	"time"

	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type SecurityBulletins struct {
	DocumentTitle      string
	ContactDetails     string
	IssuingAuthority   string
	Identification     string
	DocCreateTime      time.Time
	DocumentNotes      DocumentNotes
	AffectedComponent  string
	DocumentReferences DocumentReferences
	ProductTree        ProductTree
	Vulnerability      Vulnerability
	AffectedSystems    []dp.SystemVersion
	SeverityLevel      dp.SeverityLevel
}

func NewSecurityBulletins(defect Defect) *SecurityBulletins {
	return &SecurityBulletins{}
}

type ProductTree struct {
	Branches []Branch
}

type Branch struct {
	Type    dp.BranchType
	Name    dp.Arch
	Product []Product
}

type Product struct {
	ID       string
	CPE      string
	FullName string
}

type DocumentNotes struct {
	Synopsis    string
	Summary     string
	Description string
}

type DocumentReferences struct {
	SafetyBulletinURL dp.URL
	IssueURL          dp.URL
}

type Vulnerability struct {
	Desc        string
	ReleaseDate string
}

func (d *SecurityBulletins) SetDocumentReferences() {
	d.DocumentReferences.SafetyBulletinURL = nil
	d.DocumentReferences.IssueURL = nil
}

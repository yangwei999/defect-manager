package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type Defect struct {
	Kernel          string
	Component       string
	System          dp.SystemVersion
	ReferenceURL    dp.URL
	GuidanceURL     dp.URL
	Description     string
	SeverityLevel   dp.SeverityLevel
	AffectedSystems string
	ABI             string
	ProductTree     ProductTree
	IssueNumber     string
	IssueStatus     dp.IssueStatus
}

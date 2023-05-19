package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type Defect struct {
	Kernel          string
	Component       string
	SystemVersion   dp.SystemVersion
	ReferenceURL    dp.URL
	GuidanceURL     dp.URL
	Description     string
	SeverityLevel   dp.SeverityLevel
	AffectedVersion []dp.SystemVersion
	ABI             string
	ProductTree     ProductTree
	Issue           Issue
}

type Issue struct {
	ID     string
	Org    string
	Repo   string
	Status dp.IssueStatus
}

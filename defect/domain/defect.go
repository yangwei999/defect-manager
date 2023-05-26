package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type Defect struct {
	Kernel          string
	Component       string
	SystemVersion   dp.SystemVersion
	Description     string
	ReferenceURL    dp.URL
	GuidanceURL     dp.URL
	Influence       string
	SeverityLevel   dp.SeverityLevel
	AffectedVersion []dp.SystemVersion
	ABI             string
	Issue           Issue
}

type Issue struct {
	Number string
	Org    string
	Repo   string
	Status dp.IssueStatus
}

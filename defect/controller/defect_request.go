package controller

import (
	"github.com/opensourceways/defect-manager/defect/app"
	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

type defectRequest struct {
	IssueNumber     string   `json:"issue_number"     binding:"required"`
	IssueOrg        string   `json:"issue_org"        binding:"required"`
	IssueRepo       string   `json:"issue_repo"       binding:"required"`
	IssueStatus     string   `json:"issue_status"     binding:"required"`
	Kernel          string   `json:"kernel"           binding:"required"`
	Component       string   `json:"component"        binding:"required"`
	SystemVersion   string   `json:"system_version"   binding:"required"`
	Description     string   `json:"description"      binding:"required"`
	ReferenceURL    string   `json:"reference_url"    binding:"required"`
	GuidanceURL     string   `json:"guidance_url"     binding:"required"`
	Influence       string   `json:"influence"        binding:"required"`
	SeverityLevel   string   `json:"severity_level"   binding:"required"`
	AffectedVersion []string `json:"affected_version" binding:"required"`
	ABI             string   `json:"abi"              binding:"required"`
}

func (r defectRequest) toCmd() (cmd app.CmdToSaveDefect, err error) {
	systemVersion, err := dp.NewSystemVersion(r.SystemVersion)
	if err != nil {
		return
	}

	severityLevel, err := dp.NewSeverityLevel(r.SeverityLevel)
	if err != nil {
		return
	}

	issueStatus, err := dp.NewIssueStatus(r.IssueStatus)
	if err != nil {
		return
	}

	referenceURL, err := dp.NewURL(r.ReferenceURL)
	if err != nil {
		return
	}

	guidanceURL, err := dp.NewURL(r.GuidanceURL)
	if err != nil {
		return
	}

	affectedVersion, err := r.toAffectedVersion()
	if err != nil {
		return
	}

	cmd = app.CmdToSaveDefect{
		Kernel:          r.Kernel,
		Component:       r.Component,
		SystemVersion:   systemVersion,
		Description:     r.Description,
		ReferenceURL:    referenceURL,
		GuidanceURL:     guidanceURL,
		Influence:       r.Influence,
		SeverityLevel:   severityLevel,
		AffectedVersion: affectedVersion,
		ABI:             r.ABI,
		Issue: domain.Issue{
			Number: r.IssueNumber,
			Org:    r.IssueOrg,
			Repo:   r.IssueRepo,
			Status: issueStatus,
		},
	}

	return
}

func (r defectRequest) toAffectedVersion() (dpv []dp.SystemVersion, err error) {
	for _, v := range r.AffectedVersion {
		sv, err := dp.NewSystemVersion(v)
		if err != nil {
			return dpv, err
		}

		dpv = append(dpv, sv)
	}

	return
}

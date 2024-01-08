package app

import (
	"fmt"

	"github.com/opensourceways/defect-manager/defect/domain"
)

const (
	giteeUrl = "https://gitee.com"
)

type CmdToSaveDefect = domain.Defect

type CollectDefectsDTO struct {
	Title     string `json:"title"`
	Number    string `json:"issue_id"`
	IssueUrl  string `json:"issue_url"`
	Component string `json:"component"`
	Status    string `json:"status"`
	Score     string `json:"score"`
	Version   string `json:"version"`
}

func ToCollectDefectsDTO(defects domain.Defects) []CollectDefectsDTO {
	var dto []CollectDefectsDTO
	for _, d := range defects {
		url := fmt.Sprintf("%s/%s/%s/issues/%s", giteeUrl, d.Issue.Org, d.Issue.Repo, d.Issue.Number)

		dto = append(dto, CollectDefectsDTO{
			Title:     d.Issue.Title,
			Number:    d.Issue.Number,
			IssueUrl:  url,
			Component: d.Component,
			Status:    d.Issue.Status.String(),
			Score:     d.SeverityLevel.String(),
			Version:   d.ComponentVersion,
		})
	}

	return dto
}

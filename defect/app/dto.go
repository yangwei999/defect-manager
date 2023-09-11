package app

import "github.com/opensourceways/defect-manager/defect/domain"

type CmdToSaveDefect = domain.Defect

type CollectDefectsDTO struct {
	Number    string `json:"issue_id"`
	Component string `json:"component"`
	Status    string `json:"status"`
	Score     string `json:"score"`
	Version   string `json:"version"`
}

func ToCollectDefectsDTO(defects domain.Defects) []CollectDefectsDTO {
	var dto []CollectDefectsDTO
	for _, d := range defects {
		dto = append(dto, CollectDefectsDTO{
			Number:    d.Issue.Number,
			Component: d.Component,
			Status:    d.Issue.Status.String(),
			Score:     d.SeverityLevel.String(),
			Version:   d.ComponentVersion,
		})
	}

	return dto
}

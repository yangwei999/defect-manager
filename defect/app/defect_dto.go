package app

import (
	"time"

	"github.com/opensourceways/defect-manager/defect/domain"
)

type CmdToSaveDefect = domain.Defect

type CmdToCollectDefects struct {
	BeginTime time.Time
	Org       string
}

type CmdToGenerateBulletins struct {
	Number []string
	Org    string
}

type CollectDefectsDTO struct {
	Number    string `json:"issue_id"`
	Component string `json:"component"`
	Status    string `json:"status"`
	Score     string `json:"score"`
	Version   string `json:"version"`
}

func ToCollectDefectsDTO(defects domain.Defects) []CollectDefectsDTO {

	return nil
}

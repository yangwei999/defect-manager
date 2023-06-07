package app

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/backend"
	"github.com/opensourceways/defect-manager/defect/domain/bulletin"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
	"github.com/opensourceways/defect-manager/defect/domain/obs"
	"github.com/opensourceways/defect-manager/defect/domain/producttree"
	"github.com/opensourceways/defect-manager/defect/domain/repository"
	"github.com/opensourceways/defect-manager/utils"
)

type DefectService interface {
	SaveDefects(CmdToSaveDefect) error
	CollectDefects(time time.Time) ([]CollectDefectsDTO, error)
	GenerateBulletins([]string) error
}

func NewDefectService(
	r repository.DefectRepository,
	t producttree.ProductTree,
	b bulletin.Bulletin,
	be backend.CveBackend,
	o obs.OBS,
) *defectService {
	return &defectService{
		repo:        r,
		productTree: t,
		bulletin:    b,
		backend:     be,
		obs:         o,
	}
}

type defectService struct {
	repo        repository.DefectRepository
	productTree producttree.ProductTree
	bulletin    bulletin.Bulletin
	backend     backend.CveBackend
	obs         obs.OBS
}

func (d defectService) SaveDefects(cmd CmdToSaveDefect) error {
	has, err := d.repo.HasDefect(&cmd)
	if err != nil {
		return err
	}

	if has {
		return d.repo.SaveDefect(&cmd)
	} else {
		return d.repo.AddDefect(&cmd)
	}
}

func (d defectService) CollectDefects(date time.Time) (dto []CollectDefectsDTO, err error) {
	opt := repository.OptToFindDefects{
		BeginTime: date,
		Status:    dp.IssueStatusClosed,
	}

	defects, err := d.repo.FindDefects(opt)
	if err != nil {
		return
	}

	publishedNum, err := d.backend.IsDefectPublished(defects.AllIssueNumber())
	if err != nil {
		return
	}

	var unpublishedDefects domain.Defects
	for _, defect := range defects {
		if _, ok := publishedNum[defect.Issue.Number]; !ok {
			unpublishedDefects = append(unpublishedDefects, defect)
		}
	}

	dto = ToCollectDefectsDTO(unpublishedDefects)

	return
}

func (d defectService) GenerateBulletins(number []string) error {
	opt := repository.OptToFindDefects{
		Number: number,
	}

	defects, err := d.repo.FindDefects(opt)
	if err != nil {
		return err
	}

	maxIdentification, err := d.backend.MaxBulletinID()
	if err != nil {
		return err
	}

	bulletins := defects.GenerateBulletins()

	d.productTree.InitCache()
	defer d.productTree.CleanCache()

	for _, b := range bulletins {
		b.ProductTree, err = d.productTree.GetTree(b.Component, b.AffectedVersion)
		if err != nil {
			logrus.Errorf("component %s, get productTree error: %s", b.Component, err.Error())

			continue
		}

		maxIdentification++
		b.Identification = fmt.Sprintf("openEuler-BA-%d-%d", utils.Year(), maxIdentification)

		xmlData, err := d.bulletin.Generate(&b)
		if err != nil {
			logrus.Errorf("component: %s, to xml error: %s", b.Component, err.Error())

			continue
		}

		fileName := fmt.Sprintf("%s.xml", b.Identification)
		if err := d.obs.Upload(fileName, xmlData); err != nil {
			logrus.Errorf("component: %s, upload to obs error: %s", b.Component, err.Error())

			continue
		}
	}

	return nil
}

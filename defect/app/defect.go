package app

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/bulletin"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
	"github.com/opensourceways/defect-manager/defect/domain/producttree"
	"github.com/opensourceways/defect-manager/defect/domain/repository"
)

type DefectService interface {
	SaveDefects(CmdToSaveDefect) error
	CollectDefects(CmdToCollectDefects) ([]CollectDefectsDTO, error)
	GenerateBulletins([]string) error
}

func NewDefectService(
	r repository.DefectRepository,
	t producttree.ProductTree,
	b bulletin.Bulletin,
) *defectService {
	return &defectService{
		repo:     r,
		tree:     t,
		bulletin: b,
	}
}

type defectService struct {
	repo     repository.DefectRepository
	tree     producttree.ProductTree
	bulletin bulletin.Bulletin
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

func (d defectService) CollectDefects(cmd CmdToCollectDefects) (dto []CollectDefectsDTO, err error) {
	opt := repository.OptToFindDefects{
		BeginTime: cmd.BeginTime,
		Org:       cmd.Org,
		Status:    dp.IssueStatusClosed,
	}

	defects, err := d.repo.FindDefects(opt)
	if err != nil {
		return
	}

	return ToCollectDefectsDTO(defects), nil
}

func (d defectService) GenerateBulletins(cmd CmdToGenerateBulletins) error {
	opt := repository.OptToFindDefects{
		Org:    cmd.Org,
		Number: cmd.Number,
	}

	defects, err := d.repo.FindDefects(opt)
	if err != nil {
		return err
	}

	defectsOfComponent := defects.SeparateByComponent()
	for component, ds := range defectsOfComponent {
		var securityBulletins []domain.SecurityBulletin

		if ds.IsCombined() {
			securityBulletins = append(securityBulletins, ds.CombinedBulletin())
		} else {
			securityBulletins = append(securityBulletins, ds.SeparatedBulletins()...)
		}

		d.generateByComponent(component, securityBulletins)
	}

	return nil
}

func (d defectService) generateByComponent(component string, sbs []domain.SecurityBulletin) {
	tree, err := d.tree.GetTree(component)
	if err != nil {
		logrus.Errorf("get tree of %s err: %s", component, err.Error())

		return
	}

	for _, sb := range sbs {
		sb.Component.ProductTree = tree

		_, err := d.bulletin.Generate(sb)
		if err != nil {
			continue
		}

		//todo upload to obs
	}
}

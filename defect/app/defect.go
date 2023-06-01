package app

import (
	"fmt"

	"github.com/opensourceways/defect-manager/defect/domain/bulletin"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
	"github.com/opensourceways/defect-manager/defect/domain/producttree"
	"github.com/opensourceways/defect-manager/defect/domain/repository"
)

type DefectService interface {
	SaveDefects(CmdToSaveDefect) error
	CollectDefects(CmdToCollectDefects) ([]CollectDefectsDTO, error)
	GenerateBulletins(CmdToGenerateBulletins) error
}

func NewDefectService(
	r repository.DefectRepository,
	t producttree.ProductTree,
	b bulletin.Bulletin,
) *defectService {
	return &defectService{
		repo:        r,
		productTree: t,
		bulletin:    b,
	}
}

type defectService struct {
	repo        repository.DefectRepository
	productTree producttree.ProductTree
	bulletin    bulletin.Bulletin
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

	// todo filter published defects

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

	bulletins := defects.GenerateBulletins()

	defer d.productTree.CleanCache()
	for _, b := range bulletins {
		b.ProductTree, err = d.productTree.GetTree(b.Component, b.AffectedVersion)
		if err != nil {
			//todo log
			continue
		}

		xmlData, err := d.bulletin.Generate(&b)
		if err != nil {
			continue
		}

		fmt.Println(string(xmlData))

		//todo upload obs
	}

	return nil
}

package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
	"github.com/opensourceways/defect-manager/utils"
)

type Defects []Defect
type DefectsByComponent []Defect
type DefectsByVersion []Defect

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

func (d Defect) IsAffectVersion(version dp.SystemVersion) bool {
	for _, v := range d.AffectedVersion {
		if v == version {
			return true
		}
	}

	return false
}

func (ds Defects) SeparateByComponent() map[string]DefectsByComponent {
	classifyByComponent := make(map[string]DefectsByComponent)
	for _, d := range ds {
		classifyByComponent[d.Component] = append(classifyByComponent[d.Component], d)
	}

	return classifyByComponent
}

//IsCombined DefectsByComponent is a component-differentiated set of defects
//
func (dsc DefectsByComponent) IsCombined() bool {
	for _, d := range dsc {
		if len(d.AffectedVersion) != len(dp.MaintainVersion) {
			return false
		}

		for _, version := range d.AffectedVersion {
			if !dp.MaintainVersion[version] {
				return false
			}
		}
	}

	return true
}

//CombinedBulletin put all defect in a bulletin
func (dsc DefectsByComponent) CombinedBulletin() SecurityBulletin {
	return SecurityBulletin{
		AffectedVersion: dsc[0].AffectedVersion,
		Date:            utils.Date(),
		Component: Component{
			Name: dsc[0].Component,
		},
		Defects: Defects(dsc),
	}
}

// SeparatedBulletins separate bulletins by version name
func (dsc DefectsByComponent) SeparatedBulletins() []SecurityBulletin {
	var sbs []SecurityBulletin
	for version, ds := range dsc.separateByVersion() {
		sbs = append(sbs, ds.BulletinByVersion(version))
	}

	return sbs
}

func (dsv DefectsByVersion) BulletinByVersion(version dp.SystemVersion) SecurityBulletin {
	return SecurityBulletin{
		AffectedVersion: []dp.SystemVersion{version},
		Date:            utils.Date(),
		Component: Component{
			Name: dsv[0].Component,
		},
		Defects: Defects(dsv),
	}
}

func (dsc DefectsByComponent) separateByVersion() map[dp.SystemVersion]DefectsByVersion {
	classifyByVersion := make(map[dp.SystemVersion]DefectsByVersion)
	for version := range dp.MaintainVersion {
		for _, d := range dsc {
			if d.IsAffectVersion(version) {
				classifyByVersion[version] = append(classifyByVersion[version], d)
			}
		}
	}

	return classifyByVersion
}

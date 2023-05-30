package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
	"github.com/opensourceways/defect-manager/utils"
)

type Defects []Defect

//DefectsByComponent is group of defects by component
type DefectsByComponent []Defect

//DefectsByVersion is group of DefectsByComponent by version
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

//GroupByComponent group defects by component
func (ds Defects) GroupByComponent() map[string]DefectsByComponent {
	group := make(map[string]DefectsByComponent)
	for _, d := range ds {
		group[d.Component] = append(group[d.Component], d)
	}

	return group
}

//IsCombined DefectsByComponent is a component-differentiated set of defects,
//Bulletins are consolidated into one when all issues of a component affect all versions currently maintained,
//otherwise they are split into multiple bulletins by version
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

//CombinedBulletin put all defects in one bulletin
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

// SeparatedBulletins split into multiple bulletins by version
func (dsc DefectsByComponent) SeparatedBulletins() []SecurityBulletin {
	var sbs []SecurityBulletin
	for version, ds := range dsc.separateByVersion() {
		sbs = append(sbs, ds.bulletinByVersion(version))
	}

	return sbs
}

func (dsv DefectsByVersion) bulletinByVersion(version dp.SystemVersion) SecurityBulletin {
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

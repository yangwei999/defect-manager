package domain

import (
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

var maintainVersion = map[string]bool{
	"openEuler-20.03-LTS-SP1": true,
	"openEuler-20.03-LTS-SP3": true,
	"openEuler-22.03-LTS":     true,
	"openEuler-22.03-LTS-SP1": true,
}

type Defects []*Defect
type DefectsByComponent []*Defect

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

func (ds Defects) GenerateSecurityBulletins() []SecurityBulletin {
	classifyByComponent := make(map[string]DefectsByComponent)
	for _, d := range ds {
		classifyByComponent[d.Component] = append(classifyByComponent[d.Component], d)
	}

	// the length can't exceed len(ds)
	sbs := make([]SecurityBulletin, len(ds))

	for _, cds := range classifyByComponent {
		if cds.IsCombined() {
			sbs = append(sbs, cds.CombinedBulletin())
		} else {
			sbs = append(sbs, cds.SeparatedBulletins()...)
		}

	}

	return sbs
}

//IsCombined DefectsByComponent is a component-differentiated set of defects
//
func (cds DefectsByComponent) IsCombined() bool {
	for _, d := range cds {
		if len(d.AffectedVersion) != len(maintainVersion) {
			return false
		}

		for _, version := range d.AffectedVersion {
			if !maintainVersion[version.String()] {
				return false
			}
		}
	}

	return true
}

func (cds DefectsByComponent) CombinedBulletin() SecurityBulletin {

	return SecurityBulletin{}
}

func (cds DefectsByComponent) SeparatedBulletins() []SecurityBulletin {

	return nil
}

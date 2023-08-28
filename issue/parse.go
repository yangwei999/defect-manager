package issue

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/opensourceways/server-common-lib/utils"
	"k8s.io/apimachinery/pkg/util/sets"

	localutils "github.com/opensourceways/defect-manager/utils"
)

const (
	cmdCheck   = "/check-issue"
	cmdApprove = "/approve"

	itemKernel          = "kernel"
	itemComponents      = "components"
	itemSystemVersion   = "systemVersion"
	itemDescription     = "description"
	itemReferenceUrl    = "referenceUrl"
	itemGuidanceUrl     = "guidanceUrl"
	itemInfluence       = "influence"
	itemSeverityLevel   = "severityLevel"
	itemAffectedVersion = "affectedVersion"
	itemAbi             = "abi"

	severityLevelLow      = "Low"
	severityLevelModerate = "Moderate"
	severityLevelHigh     = "High"
	severityLevelCritical = "Critical"
)

var (
	regexpOfItems = map[string]*regexp.Regexp{
		itemKernel:          regexp.MustCompile(`(内核信息)[:：]([\s\S]*?)缺陷归属组件`),
		itemComponents:      regexp.MustCompile(`(缺陷归属组件)[:：]([\s\S]*?)组件版本`),
		itemSystemVersion:   regexp.MustCompile(`(缺陷归属的版本)[:：]([\s\S]*?)缺陷简述`),
		itemDescription:     regexp.MustCompile(`(缺陷简述)[:：]([\s\S]*?)缺陷创建时间`),
		itemReferenceUrl:    regexp.MustCompile(`(缺陷详情参考链接)[:：]([\s\S]*?)缺陷分析指导链接`),
		itemGuidanceUrl:     regexp.MustCompile(`(缺陷分析指导链接)[:：]([\s\S]*?)二、缺陷分析结构反馈`),
		itemInfluence:       regexp.MustCompile(`(影响性分析说明)[:：]([\s\S]*?)缺陷严重等级`),
		itemSeverityLevel:   regexp.MustCompile(`(缺陷严重等级)[:：]\s*(\w+)`),
		itemAffectedVersion: regexp.MustCompile(`(受影响版本排查)\(受影响/不受影响\)[:：]([\s\S]*?)abi变化`),
		itemAbi:             regexp.MustCompile(`(abi变化)\(受影响/不受影响\)[:：]([\s\S]*?)$`),
	}

	sortOfItems = []string{
		itemKernel,
		itemComponents,
		itemSystemVersion,
		itemDescription,
		itemReferenceUrl,
		itemGuidanceUrl,
		itemInfluence,
		itemSeverityLevel,
	}

	noTrimItem = map[string]bool{
		itemDescription: true,
		itemInfluence:   true,
	}

	severityLevelMap = map[string]bool{
		severityLevelLow:      true,
		severityLevelModerate: true,
		severityLevelHigh:     true,
		severityLevelCritical: true,
	}

	validCmd = []string{
		cmdCheck,
		cmdApprove,
	}
)

func (impl eventHandler) isValidCmd(comment string) bool {
	for _, v := range validCmd {
		if strings.Contains(comment, v) {
			return true
		}
	}

	return false
}

func (impl eventHandler) parseIssue(body string) (issueInfo map[string]string, err error) {
	mr := utils.NewMultiErrors()

	issueInfo = make(map[string]string)
	for _, item := range sortOfItems {
		match := regexpOfItems[item].FindAllStringSubmatch(body, -1)
		if len(match) < 1 || len(match[0]) < 3 {
			mr.Add(fmt.Sprintf("%s 解析失败", item))
			continue
		}

		trimItemInfo := localutils.TrimString(match[0][2])
		if trimItemInfo == "" {
			mr.Add(fmt.Sprintf("%s 不允许为空", match[0][1]))
			continue
		}

		if _, ok := noTrimItem[item]; ok {
			issueInfo[item] = match[0][2]
		} else {
			issueInfo[item] = trimItemInfo
		}

		switch item {
		case itemSeverityLevel:
			if _, exist := severityLevelMap[issueInfo[item]]; !exist {
				mr.Add(fmt.Sprintf("缺陷严重等级 %s 错误", issueInfo[item]))
			}

		case itemSystemVersion:
			maintainVersion := sets.NewString(impl.cfg.MaintainVersion...)
			if !maintainVersion.Has(issueInfo[item]) {
				mr.Add(fmt.Sprintf("缺陷归属版本 %s 错误", issueInfo[item]))
			}
		}
	}

	return issueInfo, mr.Err()
}

func (impl eventHandler) parseComment(s string) (affectedVersion []string, abi []string, err error) {
	mr := utils.NewMultiErrors()
	versionMatch := regexpOfItems[itemAffectedVersion].FindAllStringSubmatch(s, -1)
	if len(versionMatch) == 0 || len(versionMatch[0]) < 3 {
		mr.Add("受影响版本排查解析失败")
	}

	abiMatch := regexpOfItems[itemAbi].FindAllStringSubmatch(s, -1)
	if len(abiMatch) == 0 || len(abiMatch[0]) < 3 {
		mr.Add("abi变化解析失败")
	}
	err = mr.Err()
	if err != nil {
		return
	}

	affectedVersion, err = impl.parseVersion(versionMatch[0][2])
	if err != nil {
		return
	}

	abi, err = impl.parseVersion(abiMatch[0][2])
	if err != nil {
		return
	}

	return
}

func (impl eventHandler) parseVersion(s string) ([]string, error) {
	reg := regexp.MustCompile(`(openEuler.*?)[:：]\s*(不?受影响)`)
	matches := reg.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	var affectedVersion []string
	var allVersion []string
	for _, v := range matches {
		allVersion = append(allVersion, v[1])

		if v[2] == "受影响" {
			affectedVersion = append(affectedVersion, v[1])
		}
	}

	av := sets.NewString(allVersion...)
	if !av.HasAll(impl.cfg.MaintainVersion...) {
		return nil, fmt.Errorf("受影响版本排查/abi变化与当前维护版本不一致，当前维护版本:\n%s",
			strings.Join(impl.cfg.MaintainVersion, "\n"),
		)
	}

	return affectedVersion, nil
}

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

	itemKernel            = "kernel"
	itemComponents        = "components"
	itemComponentsVersion = "componentsVersion"
	itemSystemVersion     = "systemVersion"
	itemDescription       = "description"
	itemReferenceUrl      = "referenceUrl"
	itemGuidanceUrl       = "guidanceUrl"
	itemInfluence         = "influence"
	itemSeverityLevel     = "severityLevel"
	itemAffectedVersion   = "affectedVersion"
	itemAbi               = "abi"

	severityLevelLow      = "Low"
	severityLevelModerate = "Moderate"
	severityLevelHigh     = "High"
	severityLevelCritical = "Critical"

	regMatchResult = 0
	regMatchItem   = 2
)

var (
	itemName = map[string]string{
		itemKernel:            "内核信息",
		itemComponents:        "缺陷归属组件",
		itemComponentsVersion: "组件版本",
		itemSystemVersion:     "归属版本",
		itemDescription:       "缺陷简述",
		itemReferenceUrl:      "详情参考链接",
		itemGuidanceUrl:       "分析指导链接",
		itemInfluence:         "影响性分析说明",
		itemSeverityLevel:     "严重等级",
		itemAffectedVersion:   "受影响版本",
		itemAbi:               "abi",
	}

	regexpOfItems = map[string]*regexp.Regexp{
		itemKernel:            regexp.MustCompile(`(内核信息)[:：]([\s\S]*?)缺陷归属组件`),
		itemComponents:        regexp.MustCompile(`(缺陷归属组件)[:：]([\s\S]*?)组件版本`),
		itemComponentsVersion: regexp.MustCompile(`(组件版本)[:：]([\s\S]*?)缺陷归属的版本`),
		itemSystemVersion:     regexp.MustCompile(`(缺陷归属的版本)[:：]([\s\S]*?)缺陷简述`),
		itemDescription:       regexp.MustCompile(`(缺陷简述)[:：]([\s\S]*?)缺陷创建时间`),
		itemReferenceUrl:      regexp.MustCompile(`(缺陷详情参考链接)[:：]([\s\S]*?)缺陷分析指导链接`),
		itemGuidanceUrl:       regexp.MustCompile(`(缺陷分析指导链接)[:：]([\s\S]*?)$`),
		itemInfluence:         regexp.MustCompile(`(影响性分析说明)[:：]([\s\S]*?)缺陷严重等级`),
		itemSeverityLevel:     regexp.MustCompile(`(缺陷严重等级)[:：]([\s\S]*?)受影响版本排查`),
		itemAffectedVersion:   regexp.MustCompile(`(受影响版本排查)\(受影响/不受影响\)[:：]([\s\S]*?)abi变化`),
		itemAbi:               regexp.MustCompile(`(abi变化)\(受影响/不受影响\)[:：]([\s\S]*?)$`),
	}

	sortOfIssueItems = []string{
		itemKernel,
		itemComponents,
		itemComponentsVersion,
		itemSystemVersion,
		itemDescription,
		itemReferenceUrl,
		itemGuidanceUrl,
	}

	sortOfCommentItems = []string{
		itemInfluence,
		itemSeverityLevel,
		itemAffectedVersion,
		itemAbi,
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

type parseIssueResult struct {
	Kernel           string
	Component        string
	ComponentVersion string
	SystemVersion    string
	Description      string
	ReferenceUrl     string
	GuidanceUrl      string
}

type parseCommentResult struct {
	Influence       string
	SeverityLevel   string
	AffectedVersion []string
	Abi             []string
}

func (impl eventHandler) isValidCmd(comment string) bool {
	for _, v := range validCmd {
		if strings.Contains(comment, v) {
			return true
		}
	}

	return false
}

func (impl eventHandler) parseIssue(body string) (parseIssueResult, error) {
	result, err := impl.parse(sortOfIssueItems, body)
	if err != nil {
		return parseIssueResult{}, err
	}

	var ret parseIssueResult
	if v, ok := result[itemKernel]; ok {
		ret.Kernel = v
	}

	if v, ok := result[itemComponents]; ok {
		ret.Component = v
	}

	if v, ok := result[itemComponentsVersion]; ok {
		ret.ComponentVersion = v
	}

	if v, ok := result[itemSystemVersion]; ok {
		ret.SystemVersion = v
	}

	if v, ok := result[itemDescription]; ok {
		ret.Description = v
	}

	if v, ok := result[itemReferenceUrl]; ok {
		ret.ReferenceUrl = v
	}

	if v, ok := result[itemGuidanceUrl]; ok {
		ret.GuidanceUrl = v
	}

	return ret, nil
}

func (impl eventHandler) parseComment(body string) (parseCommentResult, error) {
	result, err := impl.parse(sortOfCommentItems, body)
	if err != nil {
		return parseCommentResult{}, err
	}

	var ret parseCommentResult
	if v, ok := result[itemInfluence]; ok {
		ret.Influence = v
	}

	if v, ok := result[itemSeverityLevel]; ok {
		ret.SeverityLevel = v
	}

	if v, ok := result[itemAffectedVersion]; ok {
		affectedVersion, err := impl.parseVersion(v)
		if err != nil {
			return parseCommentResult{}, err
		}

		ret.AffectedVersion = affectedVersion
	}

	if v, ok := result[itemAbi]; ok {
		abi, err := impl.parseVersion(v)
		if err != nil {
			return parseCommentResult{}, err
		}

		ret.Abi = abi
	}

	return ret, nil
}

func (impl eventHandler) parse(items []string, body string) (map[string]string, error) {
	mr := utils.NewMultiErrors()

	parseResult := make(map[string]string)
	for _, item := range items {
		match := regexpOfItems[item].FindAllStringSubmatch(body, -1)
		if len(match) < 1 || len(match[regMatchResult]) < 3 {
			mr.Add(fmt.Sprintf("%s 解析失败", itemName[item]))
			continue
		}

		matchItem := match[regMatchResult][regMatchItem]
		trimItemInfo := localutils.TrimString(matchItem)
		if trimItemInfo == "" {
			mr.Add(fmt.Sprintf("%s 不允许为空", itemName[item]))
			continue
		}

		if _, ok := noTrimItem[item]; ok {
			parseResult[item] = matchItem
		} else {
			parseResult[item] = trimItemInfo
		}

		switch item {
		case itemSeverityLevel:
			if _, exist := severityLevelMap[parseResult[item]]; !exist {
				mr.Add(fmt.Sprintf("缺陷严重等级 %s 错误", parseResult[item]))
			}

		case itemSystemVersion:
			maintainVersion := sets.NewString(impl.cfg.MaintainVersion...)
			if !maintainVersion.Has(parseResult[item]) {
				mr.Add(fmt.Sprintf("缺陷归属版本 %s 错误", parseResult[item]))
			}
		}
	}

	return parseResult, mr.Err()
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

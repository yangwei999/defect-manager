package issue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/defect-manager/defect/app"
	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

var Instance *eventHandler

type EventHandler interface {
	HandleIssueEvent(e *sdk.IssueEvent) error
	HandleNoteEvent(e *sdk.NoteEvent) error
}

type iClient interface {
	CreateIssueComment(org, repo string, number string, comment string) error
	ListIssueComments(org, repo, number string) ([]sdk.Note, error)
	CloseIssue(owner, repo string, number string) error
	ReopenIssue(owner, repo string, number string) error
	GetBot() (sdk.User, error)
}

func InitEventHandler(c *Config, s app.DefectService) error {
	cli := client.NewClient(func() []byte {
		return []byte(c.RobotToken)
	})

	bot, err := cli.GetBot()
	if err != nil {
		return err
	}

	Instance = &eventHandler{
		botName: bot.Login,
		cfg:     c,
		cli:     cli,
		service: s,
	}

	return nil
}

type eventHandler struct {
	botName string
	cfg     *Config
	cli     iClient
	service app.DefectService
}

func (impl eventHandler) HandleIssueEvent(e *sdk.IssueEvent) error {
	if e.Issue.TypeName != impl.cfg.IssueType {
		return nil
	}

	switch e.Issue.State {
	case sdk.StatusClosed:
		return impl.handleIssueClosed(e)

	case sdk.StatusOpen:
		return impl.handleIssueOpen(e)

	default:
		return nil
	}
}

func (impl eventHandler) handleIssueClosed(e *sdk.IssueEvent) error {
	exist, err := impl.service.IsDefectExist(&domain.Issue{
		Number: e.GetIssueNumber(),
		Org:    e.Project.Namespace,
	})
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	if err = impl.cli.ReopenIssue(e.Project.Namespace, e.Project.Name, e.Issue.Number); err != nil {
		return fmt.Errorf("reopen issue error: %s", err.Error())
	}

	logrus.Infof("reopen issue %s %s", e.Project.PathWithNamespace, e.Issue.Number)

	return impl.cli.CreateIssueComment(e.Project.Namespace,
		e.Project.Name, e.Issue.Number, "缺陷数据未收集完成，重新打开issue")
}

func (impl eventHandler) handleIssueOpen(e *sdk.IssueEvent) error {
	if _, err := impl.parseIssue(e.Issue.Body); err != nil {
		return impl.cli.CreateIssueComment(e.Project.Namespace,
			e.Project.Name, e.Issue.Number, strings.Replace(err.Error(), ". ", "\n\n", -1),
		)
	}

	return nil
}

func (impl eventHandler) HandleNoteEvent(e *sdk.NoteEvent) error {
	if !e.IsIssue() || e.Issue.TypeName != impl.cfg.IssueType ||
		e.Issue.State == sdk.StatusClosed || e.Comment.User.Login == impl.botName {
		return nil
	}

	commentIssue := func(content string) error {
		return impl.cli.CreateIssueComment(e.Project.Namespace,
			e.Project.Name, e.Issue.Number, content,
		)
	}

	if !impl.isValidCmd(e.Comment.Body) {
		if !strings.Contains(e.Comment.Body, "受影响版本排查") {
			return nil
		}

		if _, err := impl.parseComment(e.Comment.Body); err != nil {
			return commentIssue(err.Error())
		}

		return nil
	}

	issueInfo, err := impl.parseIssue(e.Issue.Body)
	if err != nil {
		return commentIssue(strings.Replace(err.Error(), ". ", "\n\n", -1))
	}

	comment := impl.approveCmdReplyToComment(e)
	if comment == "" {
		return nil
	}

	commentInfo, err := impl.parseComment(comment)
	if err != nil {
		return commentIssue(strings.Replace(err.Error(), ". ", "\n\n", -1))
	}

	if err = impl.checkRelatedPR(e, commentInfo.AffectedVersion); err != nil {
		return commentIssue(err.Error())
	}

	if err = impl.cli.CloseIssue(e.Project.Namespace, e.Project.Name, e.Issue.Number); err != nil {
		return fmt.Errorf("close issue error: %s", err.Error())
	}

	cmd, err := impl.toCmd(e, issueInfo, commentInfo)
	if err != nil {
		return fmt.Errorf("to cmd error: %s", err.Error())
	}

	err = impl.service.SaveDefects(cmd)
	if err == nil {
		return commentIssue("Your issue is accepted, thank you")
	}

	return err
}

// the content of the comment of the newest /approve reply to
func (impl eventHandler) approveCmdReplyToComment(e *sdk.NoteEvent) string {
	comments, err := impl.cli.ListIssueComments(e.Project.Namespace, e.Project.Name, e.Issue.Number)
	if err != nil {
		logrus.Errorf("get comments error: %s", err.Error())

		return ""
	}

	var id int32
	// Iterate from the end to get the latest approve command
	for i := len(comments) - 1; i >= 0; i-- {
		if strings.Contains(comments[i].Body, cmdApprove) &&
			committerInstance.isCommitter(e.Repository.PathWithNamespace, comments[i].User.Login) {
			id = comments[i].InReplyToId
			break
		}
	}
	if id == 0 {
		return ""
	}

	for _, v := range comments {
		if v.Id == id {
			return v.Body
		}
	}

	return ""
}

func (impl eventHandler) toCmd(e *sdk.NoteEvent, issue parseIssueResult, comment parseCommentResult) (
	cmd app.CmdToSaveDefect, err error) {
	systemVersion, err := dp.NewSystemVersion(issue.SystemVersion)
	if err != nil {
		return
	}

	referenceUrl, err := dp.NewURL(issue.ReferenceUrl)
	if err != nil {
		return
	}

	guidanceUrl, err := dp.NewURL(issue.GuidanceUrl)
	if err != nil {
		return
	}

	securityLevel, err := dp.NewSeverityLevel(comment.SeverityLevel)
	if err != nil {
		return
	}

	var affectedVersion []dp.SystemVersion
	for _, v := range comment.AffectedVersion {
		var dv dp.SystemVersion
		if dv, err = dp.NewSystemVersion(v); err != nil {
			return
		}
		affectedVersion = append(affectedVersion, dv)
	}

	return app.CmdToSaveDefect{
		Kernel:           issue.Kernel,
		Component:        issue.Component,
		ComponentVersion: issue.ComponentVersion,
		SystemVersion:    systemVersion,
		Description:      issue.Description,
		ReferenceURL:     referenceUrl,
		GuidanceURL:      guidanceUrl,
		Influence:        comment.Influence,
		SeverityLevel:    securityLevel,
		AffectedVersion:  affectedVersion,
		ABI:              strings.Join(comment.Abi, ","),
		Issue: domain.Issue{
			Number: e.Issue.Number,
			Org:    e.Repository.Namespace,
			Repo:   e.Repository.Name,
			Status: dp.IssueStatusClosed,
		},
	}, nil
}

func (impl eventHandler) checkRelatedPR(e *sdk.NoteEvent, versions []string) error {
	endpoint := fmt.Sprintf("https://gitee.com/api/v5/repos/%v/issues/%v/pull_requests?access_token=%s&repo=%s",
		e.Project.Namespace, e.Issue.Number, impl.cfg.RobotToken, e.Project.Name,
	)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	var prs []sdk.PullRequest
	cli := utils.NewHttpClient(3)
	bytes, _, err := cli.Download(req)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &prs); err != nil {
		return err
	}

	mergedVersion := sets.NewString()
	for _, pr := range prs {
		if pr.State == sdk.StatusMerged {
			mergedVersion.Insert(pr.Base.Ref)
		}
	}

	var relatedPRNotMerged []string
	for _, v := range versions {
		if !mergedVersion.Has(v) {
			relatedPRNotMerged = append(relatedPRNotMerged, v)
		}
	}

	if len(relatedPRNotMerged) != 0 {
		return fmt.Errorf("受影响分支关联pr未合入: %s", strings.Join(relatedPRNotMerged, ","))
	}

	return nil
}

package dp

import "errors"

const (
	open        = "open"
	progressing = "progressing"
	closed      = "closed"
	rejected    = "rejected"
)

var (
	validIssueStatus = map[string]bool{
		open:        true,
		progressing: true,
		closed:      true,
		rejected:    true,
	}

	IssueStatusClosed = issueStatus(closed)
)

type issueStatus string

type IssueStatus interface {
	String() string
}

func NewIssueStatus(s string) (IssueStatus, error) {
	if !validIssueStatus[s] {
		return nil, errors.New("invalid issue status")
	}

	return issueStatus(s), nil
}

func (s issueStatus) String() string {
	return string(s)
}

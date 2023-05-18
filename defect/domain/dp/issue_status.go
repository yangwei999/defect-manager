package dp

import "errors"

const (
	open        = "open"
	progressing = "progressing"
	closed      = "closed"
	rejected    = "rejected"
)

var validIssueStatus = map[string]bool{
	open:        true,
	progressing: true,
	closed:      true,
	rejected:    true,
}

type issueStatus string

type IssueStatus interface {
	Status() string
}

func NewIssueStatus(s string) (IssueStatus, error) {
	if !validIssueStatus[s] {
		return nil, errors.New("invalid issue status")
	}

	return issueStatus(s), nil
}

func (s issueStatus) Status() string {
	return string(s)
}

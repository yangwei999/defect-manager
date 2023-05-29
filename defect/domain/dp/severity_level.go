package dp

import "errors"

const (
	critical = "Critical"
	high     = "High"
	moderate = "Moderate"
	low      = "Low"
)

var validateSeverityLevel = map[string]bool{
	critical: true,
	high:     true,
	moderate: true,
	low:      true,
}

type severityLevel string

func NewSeverityLevel(s string) (SeverityLevel, error) {
	if !validateSeverityLevel[s] {
		return nil, errors.New("invalid severity level")
	}

	return severityLevel(s), nil
}

type SeverityLevel interface {
	String() string
}

func (s severityLevel) String() string {
	return string(s)
}

package issue

import (
	"testing"
)

func TestAssigner(t *testing.T) {
	InitCommitterInstance()

	committerInstance.initCommitterCache()
	b := committerInstance.isCommitter("src-openeuler/A-Ops", "luanjianhai")
	if !b {
		t.Failed()
	}
}

package tests

import (
	"os/exec"

	"github.com/approvals/go-approval-tests/reporters"
)

type nvimReporter struct{}

func NewNvimReporter() reporters.Reporter {
	return &nvimReporter{}
}

func (s *nvimReporter) Report(approved, received string) bool {
	_ = exec.Command("alacritty", "-e", "nvim", "-o", received, approved).Start()
	return true
}

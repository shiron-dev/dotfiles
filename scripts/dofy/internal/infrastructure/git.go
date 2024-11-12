package infrastructure

import (
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type GitInfrastructure interface {
	GitDifftool(sout io.Writer, serror io.Writer, path ...string) error
	CheckoutFile(path string) error
}

type GitInfrastructureImpl struct{}

func NewGitInfrastructure() *GitInfrastructureImpl {
	return &GitInfrastructureImpl{}
}

func (g *GitInfrastructureImpl) GitDifftool(sout io.Writer, serror io.Writer, path ...string) error {
	args := []string{"difftool", "-y"}
	args = append(args, path...)

	cmd := exec.Command("git", args...)
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "git infrastructure: failed to run git difftool")
	}

	return nil
}

func (g *GitInfrastructureImpl) CheckoutFile(path string) error {
	cmd := exec.Command("git", "checkout", "--", path)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "git infrastructure: failed to run git checkout")
	}

	return nil
}

package infrastructure

import (
	"context"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type GitInfrastructure interface {
	SetGitDir(path string)
	GitDifftool(ctx context.Context, sout io.Writer, serror io.Writer, path ...string) error
	IsGitDiff(path ...string) (bool, error)
	CheckoutFile(path string) error
}

type GitInfrastructureImpl struct {
	gitDir string
}

func NewGitInfrastructure() *GitInfrastructureImpl {
	return &GitInfrastructureImpl{
		gitDir: "",
	}
}

var ErrGitDirNotSet = errors.New("git infrastructure: git directory is not set")

func (g *GitInfrastructureImpl) SetGitDir(path string) {
	g.gitDir = path
}

func (g *GitInfrastructureImpl) GitDifftool(
	ctx context.Context,
	sout io.Writer,
	serror io.Writer,
	path ...string,
) error {
	if g.gitDir == "" {
		return ErrGitDirNotSet
	}

	args := []string{"difftool", "-y"}
	args = append(args, path...)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = g.gitDir
	cmd.Stdout = sout
	cmd.Stderr = serror

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "git infrastructure: failed to run git difftool")
	}

	return nil
}

func (g *GitInfrastructureImpl) CheckoutFile(path string) error {
	if g.gitDir == "" {
		return ErrGitDirNotSet
	}

	cmd := exec.Command("git", "checkout", "--", path)
	cmd.Dir = g.gitDir

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "git infrastructure: failed to run git checkout")
	}

	return nil
}

func (g *GitInfrastructureImpl) IsGitDiff(path ...string) (bool, error) {
	if g.gitDir == "" {
		return false, ErrGitDirNotSet
	}

	args := []string{"diff", "--quiet"}
	args = append(args, path...)

	cmd := exec.Command("git", args...)
	cmd.Dir = g.gitDir

	if err := cmd.Run(); err != nil {
		//nolint:nilerr
		return true, nil
	}

	return false, nil
}

package usecase_test

import (
	"os"
	"testing"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/test/util"
)

func TestMain(m *testing.M) {
	brew := infrastructure.NewBrewInfrastructure()

	err := util.SetupBrew(brew)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

//go:build wireinject
// +build wireinject

package di

import (
	"io"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/adapter/controller"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/usecase"

	"github.com/google/wire"
)

type stdoutType io.Writer
type stderrType io.Writer

func providePrintOutInfrastructure(stdout stdoutType, stderr stderrType) *infrastructure.PrintOutInfrastructureImpl {
	return infrastructure.NewPrintOutInfrastructure(stdout, stderr)
}

// Adapter
var controllerSet = wire.NewSet(
	wire.Bind(new(controller.DofyController), new(*controller.DofyControllerImpl)),
	controller.NewDofyController,
)

// Infrastructure
var infrastructureSet = wire.NewSet(
	wire.Bind(new(infrastructure.PrintOutInfrastructure), new(*infrastructure.PrintOutInfrastructureImpl)),
	providePrintOutInfrastructure,
	wire.Bind(new(infrastructure.ConfigInfrastructure), new(*infrastructure.ConfigInfrastructureImpl)),
	infrastructure.NewConfigInfrastructure,
	wire.Bind(new(infrastructure.BrewInfrastructure), new(*infrastructure.BrewInfrastructureImpl)),
	infrastructure.NewBrewInfrastructure,
	wire.Bind(new(infrastructure.DepsInfrastructure), new(*infrastructure.DepsInfrastructureImpl)),
	infrastructure.NewDepsInfrastructure,
	wire.Bind(new(infrastructure.FileInfrastructure), new(*infrastructure.FileInfrastructureImpl)),
	infrastructure.NewFileInfrastructure,
	wire.Bind(new(infrastructure.GitInfrastructure), new(*infrastructure.GitInfrastructureImpl)),
	infrastructure.NewGitInfrastructure,
)

// Usecase
var usecaseSet = wire.NewSet(
	wire.Bind(new(usecase.PrintOutUsecase), new(*usecase.PrintOutUsecaseImpl)),
	usecase.NewPrintOutUsecase,
	wire.Bind(new(usecase.ConfigUsecase), new(*usecase.ConfigUsecaseImpl)),
	usecase.NewConfigUsecase,
	wire.Bind(new(usecase.BrewUsecase), new(*usecase.BrewUsecaseImpl)),
	usecase.NewBrewUsecase,
	wire.Bind(new(usecase.DepsUsecase), new(*usecase.DepsUsecaseImpl)),
	usecase.NewDepsUsecase,
)

type ControllersSet struct {
	DofyController controller.DofyController
}

func InitializeControllerSet(stdout stdoutType, stderr stderrType) (*ControllersSet, error) {
	wire.Build(
		controllerSet,
		infrastructureSet,
		usecaseSet,
		wire.Struct(new(ControllersSet), "*"),
	)
	return nil, nil
}

type TestInfrastructureSet struct {
	ConfigInfrastructure   infrastructure.ConfigInfrastructure
	DepsInfrastructure     infrastructure.DepsInfrastructure
	FileInfrastructure     infrastructure.FileInfrastructure
	GitInfrastructure      infrastructure.GitInfrastructure
	PrintOutInfrastructure infrastructure.PrintOutInfrastructure
}

func InitializeTestInfrastructureSet(stdout stdoutType, stderr stderrType) (*TestInfrastructureSet, error) {
	wire.Build(
		infrastructureSet,
		wire.Struct(new(TestInfrastructureSet), "*"),
	)
	return nil, nil
}

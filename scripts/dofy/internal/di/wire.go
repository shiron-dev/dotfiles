//go:build wireinject
// +build wireinject

package di

import (
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/adapter/controller"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/usecase"

	"github.com/google/wire"
)

// Adapter
var controllerSet = wire.NewSet(
	wire.Bind(new(controller.DofyController), new(*controller.DofyControllerImpl)),
	controller.NewDofyController,
)

// Infrastructure
var infrastructureSet = wire.NewSet(
	wire.Bind(new(infrastructure.PrintOutInfrastructure), new(*infrastructure.PrintOutInfrastructureImpl)),
	infrastructure.NewPrintOutInfrastructure,
	wire.Bind(new(infrastructure.ConfigInfrastructure), new(*infrastructure.ConfigInfrastructureImpl)),
	infrastructure.NewConfigInfrastructure,
	wire.Bind(new(infrastructure.BrewInfrastructure), new(*infrastructure.BrewInfrastructureImpl)),
	infrastructure.NewBrewInfrastructure,
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

func InitializeControllerSet() (*ControllersSet, error) {
	wire.Build(
		controllerSet,
		infrastructureSet,
		usecaseSet,
		wire.Struct(new(ControllersSet), "*"),
	)
	return nil, nil
}

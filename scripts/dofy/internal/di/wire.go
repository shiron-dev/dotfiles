//go:build wireinject
// +build wireinject

package di

import (
	"dofy/internal/adapter/controller"
	"dofy/internal/infrastructure"
	"dofy/internal/usecase"

	"github.com/google/wire"
)

// Adapter
var controllerSet = wire.NewSet(
	controller.NewController,
)

// Infrastructure
var infrastructureSet = wire.NewSet(
	infrastructure.NewPrintOutInfrastructure,
)

// Usecase
var usecaseSet = wire.NewSet(
	usecase.NewPrintOutUsecase,
)

type ControllersSet struct {
	Controller controller.Controller
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

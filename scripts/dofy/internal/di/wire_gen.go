// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/adapter/controller"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/usecase"
)

// Injectors from wire.go:

func InitializeControllerSet() (*ControllersSet, error) {
	printOutInfrastructureImpl := infrastructure.NewPrintOutInfrastructure()
	printOutUsecaseImpl := usecase.NewPrintOutUsecase(printOutInfrastructureImpl)
	configInfrastructureImpl := infrastructure.NewConfigInfrastructure()
	configUsecaseImpl := usecase.NewConfigUsecase(configInfrastructureImpl)
	depsInfrastructureImpl := infrastructure.NewDepsInfrastructure()
	brewInfrastructureImpl := infrastructure.NewBrewInfrastructure()
	brewUsecaseImpl := usecase.NewBrewUsecase(brewInfrastructureImpl, printOutUsecaseImpl, configUsecaseImpl)
	depsUsecaseImpl := usecase.NewDepsUsecase(depsInfrastructureImpl, printOutUsecaseImpl, brewUsecaseImpl)
	dofyControllerImpl := controller.NewDofyController(printOutUsecaseImpl, configUsecaseImpl, depsUsecaseImpl)
	controllersSet := &ControllersSet{
		DofyController: dofyControllerImpl,
	}
	return controllersSet, nil
}

// wire.go:

// Adapter
var controllerSet = wire.NewSet(wire.Bind(new(controller.DofyController), new(*controller.DofyControllerImpl)), controller.NewDofyController)

// Infrastructure
var infrastructureSet = wire.NewSet(wire.Bind(new(infrastructure.PrintOutInfrastructure), new(*infrastructure.PrintOutInfrastructureImpl)), infrastructure.NewPrintOutInfrastructure, wire.Bind(new(infrastructure.ConfigInfrastructure), new(*infrastructure.ConfigInfrastructureImpl)), infrastructure.NewConfigInfrastructure, wire.Bind(new(infrastructure.BrewInfrastructure), new(*infrastructure.BrewInfrastructureImpl)), infrastructure.NewBrewInfrastructure, wire.Bind(new(infrastructure.DepsInfrastructure), new(*infrastructure.DepsInfrastructureImpl)), infrastructure.NewDepsInfrastructure)

// Usecase
var usecaseSet = wire.NewSet(wire.Bind(new(usecase.PrintOutUsecase), new(*usecase.PrintOutUsecaseImpl)), usecase.NewPrintOutUsecase, wire.Bind(new(usecase.ConfigUsecase), new(*usecase.ConfigUsecaseImpl)), usecase.NewConfigUsecase, wire.Bind(new(usecase.BrewUsecase), new(*usecase.BrewUsecaseImpl)), usecase.NewBrewUsecase, wire.Bind(new(usecase.DepsUsecase), new(*usecase.DepsUsecaseImpl)), usecase.NewDepsUsecase)

type ControllersSet struct {
	DofyController controller.DofyController
}

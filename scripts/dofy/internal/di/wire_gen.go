// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"dofy/internal/adapter/controller"
	"dofy/internal/infrastructure"
	"dofy/internal/usecase"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeControllerSet() (*ControllersSet, error) {
	printOutInfrastructure := infrastructure.NewPrintOutInfrastructure()
	printOutUsecase := usecase.NewPrintOutUsecase(printOutInfrastructure)
	controllerController := controller.NewController(printOutUsecase)
	controllersSet := &ControllersSet{
		Controller: controllerController,
	}
	return controllersSet, nil
}

// wire.go:

// Adapter
var controllerSet = wire.NewSet(controller.NewController)

// Infrastructure
var infrastructureSet = wire.NewSet(infrastructure.NewPrintOutInfrastructure)

// Usecase
var usecaseSet = wire.NewSet(usecase.NewPrintOutUsecase)

type ControllersSet struct {
	Controller controller.Controller
}

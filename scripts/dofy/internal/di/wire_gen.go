// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"github.com/shiron-dev/dotfiles/scripts/dofy/gen/mock/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/adapter/controller"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/usecase"
	"io"
)

// Injectors from wire.go:

func InitializeControllerSet(stdout stdoutType, stderr stderrType) (*ControllersSet, error) {
	printOutInfrastructureImpl := providePrintOutInfrastructure(stdout, stderr)
	printOutUsecaseImpl := usecase.NewPrintOutUsecase(printOutInfrastructureImpl)
	configInfrastructureImpl := infrastructure.NewConfigInfrastructure()
	configUsecaseImpl := usecase.NewConfigUsecase(configInfrastructureImpl)
	depsInfrastructureImpl := infrastructure.NewDepsInfrastructure()
	brewInfrastructureImpl := infrastructure.NewBrewInfrastructure()
	fileInfrastructureImpl := infrastructure.NewFileInfrastructure()
	gitInfrastructureImpl := infrastructure.NewGitInfrastructure()
	brewUsecaseImpl := usecase.NewBrewUsecase(brewInfrastructureImpl, depsInfrastructureImpl, printOutUsecaseImpl, configUsecaseImpl)
	depsUsecaseImpl := usecase.NewDepsUsecase(depsInfrastructureImpl, brewInfrastructureImpl, fileInfrastructureImpl, gitInfrastructureImpl, printOutUsecaseImpl, brewUsecaseImpl)
	dofyControllerImpl := controller.NewDofyController(printOutUsecaseImpl, configUsecaseImpl, depsUsecaseImpl)
	controllersSet := &ControllersSet{
		DofyController: dofyControllerImpl,
	}
	return controllersSet, nil
}

func InitializeTestInfrastructureSet(stdout stdoutType, stderr stderrType) (*TestInfrastructureSet, error) {
	brewInfrastructureImpl := infrastructure.NewBrewInfrastructure()
	configInfrastructureImpl := infrastructure.NewConfigInfrastructure()
	depsInfrastructureImpl := infrastructure.NewDepsInfrastructure()
	fileInfrastructureImpl := infrastructure.NewFileInfrastructure()
	gitInfrastructureImpl := infrastructure.NewGitInfrastructure()
	printOutInfrastructureImpl := providePrintOutInfrastructure(stdout, stderr)
	testInfrastructureSet := &TestInfrastructureSet{
		BrewInfrastructure:     brewInfrastructureImpl,
		ConfigInfrastructure:   configInfrastructureImpl,
		DepsInfrastructure:     depsInfrastructureImpl,
		FileInfrastructure:     fileInfrastructureImpl,
		GitInfrastructure:      gitInfrastructureImpl,
		PrintOutInfrastructure: printOutInfrastructureImpl,
	}
	return testInfrastructureSet, nil
}

func InitializeTestControllerSet(config *mock_infrastructure.MockConfigInfrastructure) (*TestUsecaseSet, error) {
	configUsecaseImpl := usecase.NewConfigUsecase(config)
	testUsecaseSet := &TestUsecaseSet{
		ConfigUsecase: configUsecaseImpl,
	}
	return testUsecaseSet, nil
}

// wire.go:

type (
	stdoutType io.Writer
	stderrType io.Writer
)

func providePrintOutInfrastructure(stdout stdoutType, stderr stderrType) *infrastructure.PrintOutInfrastructureImpl {
	return infrastructure.NewPrintOutInfrastructure(stdout, stderr)
}

// Adapter
var controllerSet = wire.NewSet(wire.Bind(new(controller.DofyController), new(*controller.DofyControllerImpl)), controller.NewDofyController)

// Infrastructure
var infrastructureSet = wire.NewSet(wire.Bind(new(infrastructure.PrintOutInfrastructure), new(*infrastructure.PrintOutInfrastructureImpl)), providePrintOutInfrastructure, wire.Bind(new(infrastructure.ConfigInfrastructure), new(*infrastructure.ConfigInfrastructureImpl)), infrastructure.NewConfigInfrastructure, wire.Bind(new(infrastructure.BrewInfrastructure), new(*infrastructure.BrewInfrastructureImpl)), infrastructure.NewBrewInfrastructure, wire.Bind(new(infrastructure.DepsInfrastructure), new(*infrastructure.DepsInfrastructureImpl)), infrastructure.NewDepsInfrastructure, wire.Bind(new(infrastructure.FileInfrastructure), new(*infrastructure.FileInfrastructureImpl)), infrastructure.NewFileInfrastructure, wire.Bind(new(infrastructure.GitInfrastructure), new(*infrastructure.GitInfrastructureImpl)), infrastructure.NewGitInfrastructure)

var mockInfrastructureSet = wire.NewSet()

// Usecase
var usecaseSet = wire.NewSet(wire.Bind(new(usecase.PrintOutUsecase), new(*usecase.PrintOutUsecaseImpl)), usecase.NewPrintOutUsecase, wire.Bind(new(usecase.ConfigUsecase), new(*usecase.ConfigUsecaseImpl)), usecase.NewConfigUsecase, wire.Bind(new(usecase.BrewUsecase), new(*usecase.BrewUsecaseImpl)), usecase.NewBrewUsecase, wire.Bind(new(usecase.DepsUsecase), new(*usecase.DepsUsecaseImpl)), usecase.NewDepsUsecase)

type ControllersSet struct {
	DofyController controller.DofyController
}

type TestInfrastructureSet struct {
	BrewInfrastructure     infrastructure.BrewInfrastructure
	ConfigInfrastructure   infrastructure.ConfigInfrastructure
	DepsInfrastructure     infrastructure.DepsInfrastructure
	FileInfrastructure     infrastructure.FileInfrastructure
	GitInfrastructure      infrastructure.GitInfrastructure
	PrintOutInfrastructure infrastructure.PrintOutInfrastructure
}

type TestUsecaseSet struct {
	ConfigUsecase usecase.ConfigUsecase
}

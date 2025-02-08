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
	ansibleInfrastructureImpl := infrastructure.NewAnsibleInfrastructure()
	printOutInfrastructureImpl := providePrintOutInfrastructure(stdout, stderr)
	printOutUsecaseImpl := usecase.NewPrintOutUsecase(printOutInfrastructureImpl)
	ansibleUsecaseImpl := usecase.NewAnsibleUsecase(ansibleInfrastructureImpl, printOutUsecaseImpl)
	configInfrastructureImpl := infrastructure.NewConfigInfrastructure()
	configUsecaseImpl := usecase.NewConfigUsecase(configInfrastructureImpl)
	depsInfrastructureImpl := infrastructure.NewDepsInfrastructure()
	brewInfrastructureImpl := infrastructure.NewBrewInfrastructure()
	fileInfrastructureImpl := infrastructure.NewFileInfrastructure()
	gitInfrastructureImpl := infrastructure.NewGitInfrastructure()
	brewUsecaseImpl := usecase.NewBrewUsecase(brewInfrastructureImpl, depsInfrastructureImpl, printOutUsecaseImpl, configUsecaseImpl)
	depsUsecaseImpl := usecase.NewDepsUsecase(depsInfrastructureImpl, brewInfrastructureImpl, fileInfrastructureImpl, gitInfrastructureImpl, printOutUsecaseImpl, brewUsecaseImpl)
	vsCodeInfrastructureImpl := infrastructure.NewVSCodeInfrastructure()
	vsCodeUsecaseImpl := usecase.NewVSCodeUsecase(vsCodeInfrastructureImpl, gitInfrastructureImpl, fileInfrastructureImpl, printOutUsecaseImpl, configUsecaseImpl)
	dofyControllerImpl := controller.NewDofyController(ansibleUsecaseImpl, printOutUsecaseImpl, configUsecaseImpl, depsUsecaseImpl, vsCodeUsecaseImpl)
	controllersSet := &ControllersSet{
		DofyController: dofyControllerImpl,
	}
	return controllersSet, nil
}

func InitializeTestInfrastructureSet(stdout stdoutType, stderr stderrType) (*TestInfrastructureSet, error) {
	ansibleInfrastructureImpl := infrastructure.NewAnsibleInfrastructure()
	brewInfrastructureImpl := infrastructure.NewBrewInfrastructure()
	configInfrastructureImpl := infrastructure.NewConfigInfrastructure()
	depsInfrastructureImpl := infrastructure.NewDepsInfrastructure()
	fileInfrastructureImpl := infrastructure.NewFileInfrastructure()
	gitInfrastructureImpl := infrastructure.NewGitInfrastructure()
	printOutInfrastructureImpl := providePrintOutInfrastructure(stdout, stderr)
	vsCodeInfrastructureImpl := infrastructure.NewVSCodeInfrastructure()
	testInfrastructureSet := &TestInfrastructureSet{
		AnsibleInfrastructure:  ansibleInfrastructureImpl,
		BrewInfrastructure:     brewInfrastructureImpl,
		ConfigInfrastructure:   configInfrastructureImpl,
		DepsInfrastructure:     depsInfrastructureImpl,
		FileInfrastructure:     fileInfrastructureImpl,
		GitInfrastructure:      gitInfrastructureImpl,
		PrintOutInfrastructure: printOutInfrastructureImpl,
		VSCodeInfrastructure:   vsCodeInfrastructureImpl,
	}
	return testInfrastructureSet, nil
}

func InitializeTestUsecaseSet(mockAnsibleInfrastructure *mock_infrastructure.MockAnsibleInfrastructure, mockBrewInfrastructure *mock_infrastructure.MockBrewInfrastructure, mockConfigInfrastructure *mock_infrastructure.MockConfigInfrastructure, mockDepsInfrastructure *mock_infrastructure.MockDepsInfrastructure, mockFileInfrastructure *mock_infrastructure.MockFileInfrastructure, mockGitInfrastructure *mock_infrastructure.MockGitInfrastructure, mockPrintOutInfrastructure *mock_infrastructure.MockPrintOutInfrastructure) (*TestUsecaseSet, error) {
	printOutUsecaseImpl := usecase.NewPrintOutUsecase(mockPrintOutInfrastructure)
	ansibleUsecaseImpl := usecase.NewAnsibleUsecase(mockAnsibleInfrastructure, printOutUsecaseImpl)
	configUsecaseImpl := usecase.NewConfigUsecase(mockConfigInfrastructure)
	brewUsecaseImpl := usecase.NewBrewUsecase(mockBrewInfrastructure, mockDepsInfrastructure, printOutUsecaseImpl, configUsecaseImpl)
	depsUsecaseImpl := usecase.NewDepsUsecase(mockDepsInfrastructure, mockBrewInfrastructure, mockFileInfrastructure, mockGitInfrastructure, printOutUsecaseImpl, brewUsecaseImpl)
	testUsecaseSet := &TestUsecaseSet{
		AnsibleUsecase:  ansibleUsecaseImpl,
		BrewUsecase:     brewUsecaseImpl,
		ConfigUsecase:   configUsecaseImpl,
		DepsUsecase:     depsUsecaseImpl,
		PrintOutUsecase: printOutUsecaseImpl,
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
var infrastructureSet = wire.NewSet(wire.Bind(new(infrastructure.AnsibleInfrastructure), new(*infrastructure.AnsibleInfrastructureImpl)), infrastructure.NewAnsibleInfrastructure, wire.Bind(new(infrastructure.PrintOutInfrastructure), new(*infrastructure.PrintOutInfrastructureImpl)), providePrintOutInfrastructure, wire.Bind(new(infrastructure.ConfigInfrastructure), new(*infrastructure.ConfigInfrastructureImpl)), infrastructure.NewConfigInfrastructure, wire.Bind(new(infrastructure.BrewInfrastructure), new(*infrastructure.BrewInfrastructureImpl)), infrastructure.NewBrewInfrastructure, wire.Bind(new(infrastructure.DepsInfrastructure), new(*infrastructure.DepsInfrastructureImpl)), infrastructure.NewDepsInfrastructure, wire.Bind(new(infrastructure.FileInfrastructure), new(*infrastructure.FileInfrastructureImpl)), infrastructure.NewFileInfrastructure, wire.Bind(new(infrastructure.GitInfrastructure), new(*infrastructure.GitInfrastructureImpl)), infrastructure.NewGitInfrastructure, wire.Bind(new(infrastructure.VSCodeInfrastructure), new(*infrastructure.VSCodeInfrastructureImpl)), infrastructure.NewVSCodeInfrastructure)

// Usecase
var usecaseSet = wire.NewSet(wire.Bind(new(usecase.AnsibleUsecase), new(*usecase.AnsibleUsecaseImpl)), usecase.NewAnsibleUsecase, wire.Bind(new(usecase.PrintOutUsecase), new(*usecase.PrintOutUsecaseImpl)), usecase.NewPrintOutUsecase, wire.Bind(new(usecase.ConfigUsecase), new(*usecase.ConfigUsecaseImpl)), usecase.NewConfigUsecase, wire.Bind(new(usecase.BrewUsecase), new(*usecase.BrewUsecaseImpl)), usecase.NewBrewUsecase, wire.Bind(new(usecase.DepsUsecase), new(*usecase.DepsUsecaseImpl)), usecase.NewDepsUsecase, wire.Bind(new(usecase.VSCodeUsecase), new(*usecase.VSCodeUsecaseImpl)), usecase.NewVSCodeUsecase)

type ControllersSet struct {
	DofyController controller.DofyController
}

type TestInfrastructureSet struct {
	AnsibleInfrastructure  infrastructure.AnsibleInfrastructure
	BrewInfrastructure     infrastructure.BrewInfrastructure
	ConfigInfrastructure   infrastructure.ConfigInfrastructure
	DepsInfrastructure     infrastructure.DepsInfrastructure
	FileInfrastructure     infrastructure.FileInfrastructure
	GitInfrastructure      infrastructure.GitInfrastructure
	PrintOutInfrastructure infrastructure.PrintOutInfrastructure
	VSCodeInfrastructure   infrastructure.VSCodeInfrastructure
}

type TestUsecaseSet struct {
	AnsibleUsecase  usecase.AnsibleUsecase
	BrewUsecase     usecase.BrewUsecase
	ConfigUsecase   usecase.ConfigUsecase
	DepsUsecase     usecase.DepsUsecase
	PrintOutUsecase usecase.PrintOutUsecase
}

package usecase

import (
	"github.com/pkg/errors"
	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/infrastructure"
)

type VSCodeUsecase interface {
	SaveExtensions() error
}

type VSCodeUsecaseImpl struct {
	vsCodeInfrastructure infrastructure.VSCodeInfrastructure
	gitInfrastructure    infrastructure.GitInfrastructure
	fileInfrastructure   infrastructure.FileInfrastructure
	printOutUC           PrintOutUsecase
	configUC             ConfigUsecase
}

func NewVSCodeUsecase(
	vsCodeInfrastructure infrastructure.VSCodeInfrastructure,
	gitInfrastructure infrastructure.GitInfrastructure,
	fileInfrastructure infrastructure.FileInfrastructure,
	printOutUC PrintOutUsecase,
	configUC ConfigUsecase,
) *VSCodeUsecaseImpl {
	return &VSCodeUsecaseImpl{
		vsCodeInfrastructure: vsCodeInfrastructure,
		gitInfrastructure:    gitInfrastructure,
		fileInfrastructure:   fileInfrastructure,
		printOutUC:           printOutUC,
		configUC:             configUC,
	}
}

func (v *VSCodeUsecaseImpl) SaveExtensions() error {
	v.printOutUC.PrintMdf(`
## Save VSCode extensions

`)

	dotPath, err := v.configUC.GetDotfilesDir()
	if err != nil {
		return errors.Wrap(err, "vscode usecase: failed to get dotfiles dir")
	}

	extSavePath := dotPath + "/config/vscode/extensions"

	extensions, err := v.vsCodeInfrastructure.ListExtensions()
	if err != nil {
		return errors.Wrap(err, "vscode usecase: failed to list extensions")
	}

	extensionsData := []byte{}
	for _, ext := range extensions {
		extensionsData = append(extensionsData, []byte(ext+"\n")...)
	}

	err = v.fileInfrastructure.WriteFile(extSavePath, extensionsData)
	if err != nil {
		return errors.Wrap(err, "vscode usecase: failed to write extensions")
	}

	v.gitInfrastructure.SetGitDir(dotPath)

	isDiff, err := v.gitInfrastructure.IsGitDiff(extSavePath)
	if err != nil {
		return errors.Wrap(err, "vscode usecase: failed to check git diff")
	}

	if isDiff {
		v.printOutUC.PrintMdf(`
> [!WARNING]
> **Extensions are changed.** Please check the diff and commit it.
> [!WARNING]
> **Extensions are changed.** Please check the diff and commit it.
`)

		v.printOutUC.Println(extSavePath + "\n")
	}

	return nil
}

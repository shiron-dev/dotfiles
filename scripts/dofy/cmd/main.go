package main

import (
	"os"

	"github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"
)

func main() {
	controllerSet, err := di.InitializeControllerSet(os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	dofyController := controllerSet.DofyController
	dofyController.Start()
}

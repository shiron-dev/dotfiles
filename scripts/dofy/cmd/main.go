package main

import "github.com/shiron-dev/dotfiles/scripts/dofy/internal/di"

func main() {
	controllerSet, err := di.InitializeControllerSet()
	if err != nil {
		panic(err)
	}

	dofyController := controllerSet.DofyController
	dofyController.Start()
}

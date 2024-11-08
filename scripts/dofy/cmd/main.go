package main

import "dofy/internal/di"

func main() {
	controllerSet, err := di.InitializeControllerSet()
	if err != nil {
		panic(err)
	}

	controllerSet.Controller.Start()
}

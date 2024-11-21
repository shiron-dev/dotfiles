package main

//go:generate go run github.com/google/wire/cmd/wire gen internal/di/wire.go
//go:generate go run go.uber.org/mock/mockgen -source=internal/infrastructure/brew.go -destination=./gen/mock/infrastructure/brew.go
//go:generate go run go.uber.org/mock/mockgen -source=internal/infrastructure/config.go -destination=./gen/mock/infrastructure/config.go
//go:generate go run go.uber.org/mock/mockgen -source=internal/infrastructure/deps.go -destination=./gen/mock/infrastructure/deps.go
//go:generate go run go.uber.org/mock/mockgen -source=internal/infrastructure/file.go -destination=./gen/mock/infrastructure/file.go
//go:generate go run go.uber.org/mock/mockgen -source=internal/infrastructure/git.go -destination=./gen/mock/infrastructure/git.go
//go:generate go run go.uber.org/mock/mockgen -source=internal/infrastructure/printout.go -destination=./gen/mock/infrastructure/printout.go

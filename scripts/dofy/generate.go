package main

//go:generate go run github.com/google/wire/cmd/wire gen internal/di/wire.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w -excl New.+ -template_dir internal/test/templates internal/infrastructure/brew.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w -excl New.+ -template_dir internal/test/templates internal/infrastructure/config.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w -excl New.+ -template_dir internal/test/templates internal/infrastructure/deps.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w -excl New.+ -template_dir internal/test/templates internal/infrastructure/file.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w -excl New.+ -template_dir internal/test/templates internal/infrastructure/git.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w -excl New.+ -template_dir internal/test/templates internal/infrastructure/printout.go

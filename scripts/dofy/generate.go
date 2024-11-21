package main

//go:generate go run github.com/google/wire/cmd/wire gen internal/di/wire.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w -template_dir internal/test/templates internal/infrastructure/brew.go

package main

//go:generate go run github.com/google/wire/cmd/wire gen internal/di/wire.go
//go:generate go run github.com/cweill/gotests/gotests -all -exported -parallel -w internal/infrastructure/brew.go

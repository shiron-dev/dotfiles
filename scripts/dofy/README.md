# Dofy

Dofy is a simple dotfiles manager.

## Dev dependencies

- [Go](https://go.dev//)
- [xc](https://xcfile.dev/)
- [golangci-lint](https://golangci-lint.run/)

## Tasks

> [!NOTE]
> You can use `xc`(<https://xcfile.dev/>) to run the commands
>
> See <https://xcfile.dev/getting-started/#installation> for installation instructions

### install

Install the dependencies.

Run: once

```bash
go mod tidy
```

### init

Init development environment.

Requires: install
Run: once

```bash
```

### check

Check golang code.

Requires: fmt, vet, lint
RunDeps: async

### fmt

Format golang code.

Requires: init

```bash
go fmt ./...
gofumpt -l -w .
```

### lint

Lint golang code.

Requires: init

```bash
golangci-lint run  --config=.golangci.yaml
```

### lint:fix

Lint and fix golang code.

Requires: init

```bash
golangci-lint run --fix
```

### vet

Vet golang code.

Requires: init

```bash
go vet ./...
```

### gen

Generate golang code.

Requires: init

```bash
go generate ./...
```

### test

Test golang code.

Requires: init

```bash
go test ./...
```

### build

Build the project.

Requires: init

```bash
go build -o dofy cmd/main.go
```

### clean

```bash
go clean -testcache
```

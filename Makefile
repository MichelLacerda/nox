# Variáveis para configuração de build
OS ?= windows
ARCH ?= amd64
BIN_EXT = $(if $(filter $(OS),windows),.exe,)

.PHONY: all
all: test build build-fmt

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "-w -s" -o ./bin/nox$(BIN_EXT) ./cmd/nox/main.go

.PHONY: build-fmt
build-fmt:
	GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "-w -s" -o ./bin/noxfmt$(BIN_EXT) ./cmd/fmt/main.go

.PHONY: clean
clean:
	rm -f ./bin/nox$(BIN_EXT)
	rm -f ./bin/noxfmt$(BIN_EXT)

.PHONY: run
run: build
	./bin/nox$(BIN_EXT)

.PHONY: install
install: build build-fmt
	@if [ "$(OS)" = "windows" ]; then \
		mkdir "%USERPROFILE%\\.local\\bin" || exit 0; \
		copy ".\\bin\\nox$(BIN_EXT)" "%USERPROFILE%\\.local\\bin\\nox$(BIN_EXT)"; \
		copy ".\\bin\\noxfmt$(BIN_EXT)" "%USERPROFILE%\\.local\\bin\\noxfmt$(BIN_EXT)"; \
	else \
		mkdir -p $(HOME)/.local/bin || exit 0; \
		cp ./bin/nox$(BIN_EXT) $(HOME)/.local/bin/nox$(BIN_EXT); \
		cp ./bin/noxfmt$(BIN_EXT) $(HOME)/.local/bin/noxfmt$(BIN_EXT); \
	fi

.PHONY: test-examples
test-examples-w64:
	powershell -ExecutionPolicy Bypass -File .\examples\test-examples.ps1

.PHONY: test-examples-linux
test-examples-linux:
	sh ./examples/test-examples.sh
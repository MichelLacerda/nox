GOOS := $(shell go env GOOS)
BIN_EXT = $(if $(filter $(GOOS),windows),.exe,)

.PHONY: all
all: test build build-fmt

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build -ldflags "-w -s" -o ./bin/nox$(BIN_EXT) ./cmd/nox/main.go

.PHONY: build-fmt
build-fmt:
	go build -ldflags '-w -s' -o ./bin/noxfmt$(BIN_EXT) ./cmd/fmt/main.go

.PHONY: clean
clean:
	rm -f ./bin/nox$(BIN_EXT)
	rm -f ./bin/noxfmt$(BIN_EXT)

.PHONY: run
run: build
	./bin/nox$(BIN_EXT)

.PHONY: install-windows
install-windows: build build-fmt
	mkdir "%USERPROFILE%\\.local\\bin" || exit 0
	copy ".\\bin\\nox.exe" "%USERPROFILE%\\.local\\bin\\nox.exe"
	copy ".\\bin\\noxfmt.exe" "%USERPROFILE%\\.local\\bin\\noxfmt.exe"

.PHONY: install-linux
install-linux: build build-fmt
	mkdir -p $(HOME)/.local/bin || exit 0; \
	cp ./bin/nox$(BIN_EXT) $(HOME)/.local/bin/nox$(BIN_EXT); \
	cp ./bin/noxfmt$(BIN_EXT) $(HOME)/.local/bin/noxfmt$(BIN_EXT); \

.PHONY: test-examples-windows
test-examples-windows:
	powershell -ExecutionPolicy Bypass -File .\examples\test-examples.ps1

.PHONY: test-examples-linux
test-examples-linux:
	sh ./examples/test-examples.sh
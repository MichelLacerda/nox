GOOS := $(shell go env GOOS)
BIN_EXT = $(if $(filter $(GOOS),windows),.exe,)

.PHONY: all
all: test build build-fmt

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build -ldflags "-w -s" -o ./bin/$(GOOS)/nox$(BIN_EXT) -trimpath ./cmd/nox/main.go

.PHONY: build-fmt
build-fmt:
	go build -ldflags '-w -s' -o ./bin/$(GOOS)/noxfmt$(BIN_EXT) -trimpath ./cmd/fmt/main.go

.PHONY: clean
clean:
	rm -f ./bin/$(GOOS)/nox$(BIN_EXT)
	rm -f ./bin/$(GOOS)/noxfmt$(BIN_EXT)

.PHONY: run
run: build
	./bin/$(GOOS)/nox$(BIN_EXT)

.PHONY: install-windows
install-windows: build build-fmt
	mkdir "%USERPROFILE%\\.local\\bin" || exit 0
	copy ".\\bin\\$(GOOS)\\nox.exe" "%USERPROFILE%\\.local\\bin\\nox.exe"
	copy ".\\bin\\$(GOOS)\\noxfmt.exe" "%USERPROFILE%\\.local\\bin\\noxfmt.exe"

.PHONY: install-linux
install-linux: build build-fmt
	mkdir -p $(HOME)/.local/bin || exit 0; \
	cp ./bin/$(GOOS)/nox$(BIN_EXT) $(HOME)/.local/bin/nox$(BIN_EXT); \
	cp ./bin/$(GOOS)/noxfmt$(BIN_EXT) $(HOME)/.local/bin/noxfmt$(BIN_EXT); \

.PHONY: test-examples-windows
test-examples-windows:
	powershell -ExecutionPolicy Bypass -File .\examples\test-examples.ps1

.PHONY: test-examples-linux
test-examples-linux:
	sh ./examples/test-examples.sh
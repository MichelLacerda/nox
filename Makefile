.PHONY: all
all: test build-w64 build-fmt-w64

.PHONY: build-w64
build-w64:
	go build -ldflags "-w -s" -o ./bin/nox.exe .\cmd\nox\main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: clean-w64
clean-w64:
	rm ./bin/nox.exe
	rm ./bin/noxfmt.exe

.PHONY: run-w64
run-w64: build-w64
	./bin/nox.exe

.PHONY: stats
stats:
	cloc . --exclude-dir=vendor,assets,docs,examples,tests

.PHONY: build-fmt-w64
build-fmt-w64:
	go build -ldflags "-w -s" -o ./bin/noxfmt.exe .\cmd\fmt\main.go

.PHONY: install-w64
install-w64: build-w64 build-fmt-w64
	mkdir "%USERPROFILE%\\.local\\bin" || exit 0
	copy ".\\bin\\nox.exe" "%USERPROFILE%\\.local\\bin\\nox.exe"
	copy ".\\bin\\noxfmt.exe" "%USERPROFILE%\\.local\\bin\\noxfmt.exe"

.PHONY: test-examples-w64
test-examples-w64:
	powershell -ExecutionPolicy Bypass -File .\examples\test-examples.ps1

.PHONY: test-examples-linux
test-examples-linux:
	@find examples -type f -name "*.nox" | while read file; do \
		echo "Executing script: $$(basename $$file)"; \
		./nox "$$file"; \
	done
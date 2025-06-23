.PHONY: all
all: test build build-fmt

.PHONY: build
build:
	go build -ldflags "-w -s" -o ./bin/nox.exe .\cmd\nox\main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm ./bin/nox.exe
	rm ./bin/noxfmt.exe

.PHONY: run
run: build
	./bin/nox.exe

.PHONY: stats
stats:
	cloc . --exclude-dir=vendor,assets,docs,examples,tests

.PHONY: build-fmt
build-fmt:
	go build -ldflags "-w -s" -o ./bin/noxfmt.exe .\cmd\fmt\main.go

.PHONY: install
install: build build-fmt
	mkdir "%USERPROFILE%\\.local\\bin" || exit 0
	copy ".\\bin\\nox.exe" "%USERPROFILE%\\.local\\bin\\nox.exe"
	copy ".\\bin\\noxfmt.exe" "%USERPROFILE%\\.local\\bin\\noxfmt.exe"
.PHONY: all
all: test build build-fmt

.PHONY: build
build:
	go build -ldflags "-w -s" -o ./bin/nox.exe .

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -f ./bin/nox.exe
	rm -f ./bin/noxfmt.exe

.PHONY: run
run: build
	./nox.exe

.PHONY: stats
stats:
	cloc . --exclude-dir=vendor,assets,docs,examples,tests

.PHONY: build-fmt
build-fmt:
	go build -o ./bin/noxfmt.exe .\cmd\fmt\main.go


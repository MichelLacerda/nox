.PHONY: build

build:
	go build -ldflags "-w -s" .

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	go clean
	rm -f ./nox.exe

.PHONY: run
run: build
	./nox.exe

.PHONY: stats
stats:
	cloc . --exclude-dir=vendor,assets,docs,examples,tests
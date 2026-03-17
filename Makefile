.PHONY: build build-all clean test install

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/terminal-help/th/cmd.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o th .

build-all: clean
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/th-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/th-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/th-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/th-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/th-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/th-windows-arm64.exe .

clean:
	rm -rf dist/
	rm -f th

test:
	go test -v ./...

install: build
	cp th /usr/local/bin/th

dev:
	go run .

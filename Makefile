GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
VERSION=0.1.5
BINARY_NAME=chatgpt-wecom
LDFLAGS=-ldflags '-linkmode external -extldflags "-static"'

all: mac-amd64 mac-arm64 linux-amd64 linux-arm64 win64

dockerenv:
	 docker build -t ${BINARY_NAME}:${VERSION} -f $(shell pwd)/docker/callback.Dockerfile .

mac-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -ldflags "-s -w " -o $(BINARY_NAME).$(VERSION).amd64-darwin ./cmd/app

mac-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GOBUILD) -ldflags "-s -w " -o $(BINARY_NAME).$(VERSION).arm64-darwin ./cmd/app

linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME).$(VERSION).amd64-linux ./cmd/app

linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-unknown-linux-gnu-gcc $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME).$(VERSION).arm64-linux ./cmd/app

win64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME).$(VERSION).exe ./cmd/app

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME).$(VERSION).amd64-linux $(BINARY_NAME).$(VERSION).amd64-darwin $(BINARY_NAME).$(VERSION).arm64-darwin $(BINARY_NAME).$(VERSION).arm64-linux $(BINARY_NAME).$(VERSION).exe
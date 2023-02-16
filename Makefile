GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
VERSION=0.1.0
BINARY_NAME=chatgpt-wecom

all: amd64 arm64 win64

amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-s -w " -o $(BINARY_NAME).$(VERSION).amd64 ./cmd/app

arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags "-s -w " -o $(BINARY_NAME).$(VERSION).arm64 ./cmd/app

win64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "-s -w -H windowsgui" -o $(BINARY_NAME).$(VERSION).exe ./cmd/app

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME).$(VERSION).amd64 $(BINARY_NAME).$(VERSION).arm64 $(BINARY_NAME).$(VERSION).exe
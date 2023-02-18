GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
VERSION=0.1.1
BINARY_NAME=chatgpt-wecom

all: amd64 arm64 win64 mac

mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "-s -w " -o $(BINARY_NAME).$(VERSION).amd64-darwin ./cmd/app

amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-s -w " -o $(BINARY_NAME).$(VERSION).amd64-linux ./cmd/app

arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags "-s -w " -o $(BINARY_NAME).$(VERSION).arm64-linux ./cmd/app

win64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "-s -w -H windowsgui" -o $(BINARY_NAME).$(VERSION).exe ./cmd/app

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME).$(VERSION).amd64-linux $(BINARY_NAME).$(VERSION).amd64-darwin $(BINARY_NAME).$(VERSION).arm64-linux $(BINARY_NAME).$(VERSION).exe
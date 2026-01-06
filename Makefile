projectName := file-uploader
VERSION = 0.0.1

# Cross-compilation values.
MAC_ARCH=arm64
LINUX_ARCH=amd64
OS_LINUX=linux
OS_MAC=linux

# Output directory structures.
BUILD=build
LINUX_BUILD_ARCH=./$(OS_LINUX)-$(LINUX_ARCH)
MAC_BUILD_ARCH=./$(OS_MAC)-$(MAC_ARCH)

# Flags to pass to the Go linker using the -ldflags="-X ..." option.
PACKAGE_PATH=github.com/ONSdigital/ras-rm-sample-file-uploader
BRANCH_FLAG=$(PACKAGE_PATH)/models.branch=$(BRANCH)
BUILT_FLAG=$(PACKAGE_PATH)/models.built=$(BUILT)
COMMIT_FLAG=$(PACKAGE_PATH)/models.commit=$(COMMIT)
ORIGIN_FLAG=$(PACKAGE_PATH)/models.origin=$(ORIGIN)
VERSION_FLAG=$(PACKAGE_PATH)/models.version=$(VERSION)

# Get the Git branch the commit is from, stripping the leading asterisk.
export BRANCH?=$(shell git branch --contains $(COMMIT) | grep \* | cut -d ' ' -f2)

# Get the current date/time in UTC and ISO-8601 format.
export BUILT?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Get the full Git commit SHA-1 hash.
export COMMIT?=$(shell git rev-parse HEAD)

# Get the Git repo origin.
export ORIGIN?=$(shell git config --get remote.origin.url)

# Cross-compile the binary for Linux and macOS, setting linker flags for information returned by the GET /info endpoint.
build: clean
	GOOS=$(OS_LINUX) GOARCH=$(LINUX_ARCH) CGO_ENABLED=0 go build -o $(LINUX_BUILD_ARCH)-main -ldflags="-X $(BRANCH_FLAG) -X $(BUILT_FLAG) -X $(COMMIT_FLAG) -X $(ORIGIN_FLAG) -X $(VERSION_FLAG)" main.go
	GOOS=$(OS_MAC) GOARCH=$(MAC_ARCH) CGO_ENABLED=0 go build -o $(MAC_BUILD_ARCH)-main -ldflags="-X $(BRANCH_FLAG) -X $(BUILT_FLAG) -X $(COMMIT_FLAG) -X $(ORIGIN_FLAG) -X $(VERSION_FLAG)" main.go

# Remove the build directory tree.
clean:
	if [ -d $(BUILD) ]; then rm -r linux-*-main; fi;

.PHONY: install
install: checkgo clean mod
	go clean -i && go build -v -o main
	chmod 755 main

.PHONY: checkgo
checkgo:
ifndef GOPATH
	$(error GOPATH is not set)
endif

.PHONY: mod
mod:
	go mod download

.PHONY: test
test:
	go test -v -parallel=1 -race -coverprofile=coverage.txt ./...

projectName := file-uploader

.PHONY: install
install: checkgo clean mod
	go clean -i && go build -v -o main
	chmod 755 main

.PHONY: checkgo
checkgo:
ifndef GOPATH
	$(error GOPATH is not set)
endif

.PHONY: clean
clean:
	rm -f main

.PHONY: mod
mod:
	go mod download

# TODO: Remove this once a fix for https://github.com/golang/go/issues/75031 is released,
# added in https://github.com/ONSdigital/ras-rm-sample-file-uploader/pull/51.
# This is a workaround for an issue where `go test` fails in GHA because it cannot find the
# `covdata` package caused by the above issue with the `GOTOOLCHAIN` version
.PHONY: test
test:
	GOVERSION := $(shell go env GOVERSION)
	export GOTOOLCHAIN := $(GOVERSION)+auto
	go test -v -parallel=1 -race -coverprofile=coverage.txt ./...

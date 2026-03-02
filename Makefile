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

.PHONY: test
test:
	go test -v -parallel=1 -race -coverprofile=coverage.txt ./...

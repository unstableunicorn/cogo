OUT := cogo
PKG := github.com/unstableunicorn/cogo
VERSION := $(shell git describe --always --long --dirty)
VERSION_PATH := github.com/unstableunicorn/cogo/cmd.version
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)


all: run

buildlocal: 
	go build -i -v -o bin/${OUT} -ldflags="-X ${VERSION_PATH}=${VERSION}" ${PKG}

build: static
	GOOS=darwin GOARCH=amd64 go build -i -v -o bin/${OUT}-darwin-amd64 -ldflags="-X ${VERSION_PATH}=${VERSION}" ${PKG}
	GOOS=linux GOARCH=amd64 go build -i -v -o bin/${OUT}-linux-amd64 -ldflags="-X ${VERSION_PATH}=${VERSION}" ${PKG}
	GOOS=windows GOARCH=amd64 go build -i -v -o bin/${OUT}-windows-amd64 -ldflags="-X ${VERSION_PATH}=${VERSION}" ${PKG}

test:
	@go test -short ${PKG_LIST}

getmodules:
	go mod download

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint -set_exit_status $$file ; \
	done

static: vet lint
	go build -i -v -o bin/${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X cmd.version=${VERSION}" ${PKG}

run: buildlocal
	./${OUT}

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run server static vet lint
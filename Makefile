PROJECT_NAME=mail-checker
BUILD_VERSION=$(shell cat VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on
GO_BUILD_WINDOWS_ENV=CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on
GO_FILES=$(shell go list ./... | grep -v /vendor/)
GO_VERSION=1.21.6

.SILENT:

all: mod_tidy fmt vet  install test build

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin .
	$(GO_BUILD_WINDOWS_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).exe .

install:
	$(GO_BUILD_ENV) go install

vet:
	$(GO_BUILD_ENV) go vet $(GO_FILES)

fmt:
	$(GO_BUILD_ENV) go fmt $(GO_FILES)

mod_tidy:
	$(GO_BUILD_ENV) go mod tidy -compat=$(GO_VERSION)

test:
	go test -cover
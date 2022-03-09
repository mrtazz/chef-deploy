#
# some housekeeping tasks
#

VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)

GOFLAGS :=
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.goversion=$(GOVERSION)'

.PHONY: build
build:
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" cmd/chef-deploy.go

.PHONY: test
test:
	go test $(GOFLAGS) -v ./...

.PHONY: coverage
coverage:
	go test $(GOFLAGS) -coverprofile=cover.out ./...
	go tool $(GOFLAGS) cover -html=cover.out -o cover.html

.PHONY: benchmark
benchmark:
	@echo "Running tests..."
	@go test $(GOFLAGS) -bench=. ${NAME}

.DEFAULT_GOAL:=build

#
# some housekeeping tasks
#

NAME := chef-deploy
DESC := tool to deploy changes to a Chef server based on a git diff
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)


BUILD_GOOS ?= $(shell go env GOOS)
BUILD_GOARCH ?= $(shell go env GOARCH)

RELEASE_ARTIFACTS_DIR := .release_artifacts
CHECKSUM_FILE := checksums.txt

GOFLAGS :=
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.goversion=$(GOVERSION)'

.PHONY: build
build: chef-deploy

.PHONY: chef-deploy
chef-deploy: chef-deploy.go
	GOOS=$(BUILD_GOOS) GOARCH=$(BUILD_GOARCH) go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $@ $<
.DEFAULT_GOAL:=build

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

.PHONY: build-artifact
build-artifact: certcal $(RELEASE_ARTIFACTS_DIR)
	mv certcal $(RELEASE_ARTIFACTS_DIR)/certcal-$(VERSION).$(BUILD_GOOS).$(BUILD_GOARCH)
	cd $(RELEASE_ARTIFACTS_DIR) && shasum -a 256 certcal-$(VERSION).$(BUILD_GOOS).$(BUILD_GOARCH) >> $(CHECKSUM_FILE)

.PHONY: github-release
github-release:
	gh release create $(VERSION) --title 'Release $(VERSION)' --notes-file docs/releases/$(VERSION).md $(RELEASE_ARTIFACTS_DIR)/*

# clean up tasks
.PHONY: clean
clean:
	git clean -fdx

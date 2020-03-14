#
# some housekeeping tasks
#

VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)

GOFLAGS := -mod=vendor
LDFLAGS := -X 'github.com/mrtazz/chef-deploy/pkg/version.version=$(VERSION)' \
           -X 'github.com/mrtazz/chef-deploy/pkg/version.goversion=$(GOVERSION)'

build:
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" cmd/chef-deploy.go

vendor:
	go mod vendor

test:
	go test $(GOFLAGS) -v $$(go list $(GOFLAGS) ./... | grep -v /cmd/)

coverage:
	@echo "mode: set" > cover.out
	@for package in $(PACKAGES); do \
		if ls $${package}/*.go &> /dev/null; then  \
		go test $(GOFLAGS) -coverprofile=$${package}/profile.out $${package}; fi; \
		if test -f $${package}/profile.out; then \
		cat $${package}/profile.out | grep -v "mode: set" >> cover.out; fi; \
	done
	@-go tool $(GOFLAGS) cover -html=cover.out -o cover.html

benchmark:
	@echo "Running tests..."
	@go test $(GOFLAGS) -bench=. ${NAME}

.DEFAULT_GOAL:=build

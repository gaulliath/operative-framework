# Build targets

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')
PACKAGES = $(shell go list ./...)

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION_GIT := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
VERSION := $(if $(VERSION),$(VERSION),$(VERSION_GIT))
BIND_DIR := "dist"
GIT_BRANCH := $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))
REPONAME := $(shell echo $(REPO) | tr '[:upper:]' '[:lower:]')

build-docker:
	@docker build -t opf --no-cache .

build-docker-compose:
	@docker-compose build

## Format the Code
fmt:
	gofmt -s -l -w $(SRCS)

check: vet lint errcheck interfacer aligncheck structcheck varcheck unconvert staticcheck vendorcheck prealloc test

vet:
	go vet $(PACKAGES)

lint:
	golint -set_exit_status $(PACKAGES)

errcheck:
	errcheck -exclude errcheck_excludes.txt $(PACKAGES)

interfacer:
	interfacer $(PACKAGES)

aligncheck:
	aligncheck $(PACKAGES)

structcheck:
	structcheck $(PACKAGES)

varcheck:
	varcheck $(PACKAGES)

unconvert:
	unconvert -v $(PACKAGES)

staticcheck:
	staticcheck $(PACKAGES)

vendorcheck:
	vendorcheck $(PACKAGES)
	vendorcheck -u $(PACKAGES)

prealloc:
	prealloc $(PACKAGES)

test:
	go test -cover $(PACKAGES)

coverage:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

deps:
	go get -u github.com/alexkohler/prealloc
	go get -u github.com/FiloSottile/vendorcheck
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/mdempsky/unconvert
	go get -u github.com/opennota/check/...
	go get -u honnef.co/go/tools/...
	go get -u mvdan.cc/interfacer
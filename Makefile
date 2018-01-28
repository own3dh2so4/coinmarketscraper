export

# Project specific variables
PROJECT=coinmarketscraper
OS=$(shell uname)
VENDOR=./src/vendor

# GO env
GOPATH=$(shell pwd)
GO=go
GOCMD=GOPATH=$(GOPATH) $(GO)

# Build versioning
COMMIT = $(shell git log -1 --format="%h" 2>/dev/null || echo "0")
VERSION=$(shell git describe --tags --always)
BUILD_DATE = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
FLAGS = -ldflags "\
  -X $(PROJECT)/constants.COMMIT=$(COMMIT) \
  -X $(PROJECT)/constants.VERSION=$(VERSION) \
  -X $(PROJECT)/constants.BUILD_DATE=$(BUILD_DATE) \
  "

GOBUILD = $(GOCMD) build $(FLAGS)

.PHONY: all
all:	build


.PHONY: build
build: format test compile

.PHONY: compile
compile:
	$(GOBUILD) -o bin/$(PROJECT) $(PROJECT)

.PHONY: format
format:
	@for gofile in $$(find ./src/$(PROJECT) -name "*.go"); do \
		echo "formatting" $$gofile; \
		gofmt -w $$gofile; \
	done

.PHONY: goreport
goreport:
	@for gofile in $$(find ./src/$(PROJECT) -name "*.go"); do \
		echo "cmd: gofmt -w " $$gofile; \
		gofmt -w $$gofile; \
		echo "cmd: go tool vet " $$gofile; \
		go tool vet $$gofile; \
		echo "cmd: golint " $$gofile; \
		golint $$gofile; \
		echo "cmd: gocyclo -over 15 " $$gofile; \
		gocyclo -over 15 $$gofile; \
		echo "cmd: ineffassign " $$gofile; \
		ineffassign $$gofile; \
		echo "cmd: misspell " $$gofile; \
		misspell $$gofile; \
	done


.PHONY: run
run: build
	$(GOPATH)/bin/$(PROJECT)

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: metrics
metrics:
	gometalinter.v2 --checkstyle src/$(PROJECT)/... > report.xml
	sed -i 's#src/titan/src/titan#src/titan#g' report.xml
	@for D in $$(find ./src/$(PROJECT) -type d ); do \
		go test $$D -coverprofile=cover.out; \
		gocov convert cover.out | gocov-xml > $$D/coverage.xml; \
	done
	go test -v ./src/$(PROJECT)/... | go-junit-report > test.xml

.PHONY: test
test:
	$(GOCMD) test -v -race ./src/$(PROJECT)/... -cover

.PHONY: coverage
coverage:
	rm -fr coverage
	mkdir -p coverage
	$(GOCMD) list $(PROJECT)/... > coverage/packages
	@i=a ; \
	while read -r P; do \
		i=a$$i ; \
		$(GOCMD) test ./src/$$P -cover -covermode=count -coverprofile=coverage/$$i.out; \
	done <coverage/packages

	echo "mode: count" > coverage/coverage
	cat coverage/*.out | grep -v "mode: count" >> coverage/coverage
	$(GOCMD) tool cover -html=coverage/coverage

.PHONY: cleanvendor
cleanvendor:
	@for pattern in *_test.go .travis.yml LICENSE Makefile CONTRIBUTORS .gitattributes AUTHORS PATENTS README .gitignore; do \
		echo 'Deleting' $$pattern ; \
		find src/vendor -name $$pattern -delete; \
	done
	@for pattern in */testdata/* */.git/*; do \
		echo 'Deleting' $$pattern ; \
		find src/vendor -path $$pattern -delete; \
	done

multi: build armv5 armv6 armv7 armv8 darwin linux mipsel


armv5:
	GOOS=linux GOARM=5 GOARCH=arm $(GOBUILD) -o bin/$(PROJECT)_armv5 $(PROJECT)
armv6:
	GOOS=linux GOARM=6 GOARCH=arm $(GOBUILD) -o bin/$(PROJECT)_armv6 $(PROJECT)
armv7:
	GOOS=linux GOARM=7 GOARCH=arm $(GOBUILD) -o bin/$(PROJECT)_armv7 $(PROJECT)
armv8:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o bin/$(PROJECT)_armv8 $(PROJECT)
darwin:
	GOOS=darwin $(GOBUILD) -o bin/$(PROJECT)_darwin $(PROJECT)
linux:
	GOOS=linux $(GOBUILD) -o bin/$(PROJECT)_linux $(PROJECT)
mipsel:
	GOOS=linux GOARCH=mipsle $(GOBUILD) -o bin/$(PROJECT)_mipsel $(PROJECT)
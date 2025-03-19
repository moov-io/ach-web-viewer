# generated-from:0f1cfb3f9faa0c83355794c5720cb80c30b77f4fcb2887d31d2887bd169db413 DO NOT REMOVE, DO UPDATE

PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
PWD := $(shell pwd)

ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

COMMIT_HASH :=$(shell git rev-parse --short HEAD)
DEV_VERSION := dev-${COMMIT_HASH}

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

export GOPRIVATE=github.com/moov-io

all: install build

.PHONY: install
install:
	go mod tidy

build:
	go build -ldflags "-X github.com/moov-io/ach-web-viewer.Version=${VERSION}" -o bin/ach-web-viewer github.com/moov-io/ach-web-viewer/cmd/ach-web-viewer

.PHONY: setup
setup:
	docker compose up -d --force-recreate --remove-orphans

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=10.0 ./lint-project.sh
endif

.PHONY: teardown
teardown:
	-docker compose down --remove-orphans

docker:
	docker build --pull --build-arg VERSION=${VERSION} -t moov/ach-web-viewer:${VERSION} -f Dockerfile .

docker-push:
	docker push moov/ach-web-viewer:${VERSION}

.PHONY: dev-docker
dev-docker:
	docker build --pull --build-arg VERSION=${DEV_VERSION} -t moov/ach-web-viewer:${DEV_VERSION} -f Dockerfile .

.PHONY: dev-push
dev-push:
	docker push moov/ach-web-viewer:${DEV_VERSION}

# Extra utilities not needed for building

run: build
	./bin/ach-web-viewer

docker-run:
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm moov/ach-web-viewer:${VERSION}

test:
	go test -cover github.com/moov-io/ach-web-viewer/...

.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@echo "Skipping cleanup on Windows, currently unsupported."
else
	@rm -rf cover.out coverage.txt misspell* staticcheck*
	@rm -rf ./bin/
endif

# For open source projects

dist: clean build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/ach-web-viewer.exe cmd/ach-web-viewer/*
else
	CGO_ENABLED=1 GOOS=$(PLATFORM) go build -o bin/ach-web-viewer-$(PLATFORM)-amd64 cmd/ach-web-viewer/*
endif

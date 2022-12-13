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

all: install update build

.PHONY: install
install:
	go mod tidy
	go install github.com/markbates/pkger/cmd/pkger@latest
	go mod vendor

update-pkger:
	pkger -include /configs/config.default.yml -include /webui

update: update-pkger
	go mod vendor

build:
	go build -mod=vendor -ldflags "-X github.com/moov-io/ach-web-viewer.Version=${VERSION}" -o bin/ach-web-viewer github.com/moov-io/ach-web-viewer/cmd/ach-web-viewer

.PHONY: setup
setup:
	docker-compose up -d --force-recreate --remove-orphans

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	DISABLE_GITLEAKS=true ./lint-project.sh
endif

.PHONY: teardown
teardown:
	-docker-compose down --remove-orphans

docker: update
	docker build --pull --build-arg VERSION=${VERSION} -t moov-io/ach-web-viewer:${VERSION} -f Dockerfile .
	docker tag moov-io/ach-web-viewer:${VERSION} moov-io/ach-web-viewer:latest

	docker tag moov-io/ach-web-viewer:${VERSION} moov/ach-web-viewer:${VERSION}
	docker tag moov-io/ach-web-viewer:${VERSION} moov/ach-web-viewer:latest


docker-push:
	docker push moov/ach-web-viewer:${VERSION}
	docker push moov/ach-web-viewer:latest

.PHONY: dev-docker
dev-docker: update
	docker build --pull --build-arg VERSION=${DEV_VERSION} -t moov-io/ach-web-viewer:${DEV_VERSION} -f Dockerfile .
	docker tag moov-io/ach-web-viewer:${DEV_VERSION} moov/ach-web-viewer:${DEV_VERSION}

.PHONY: dev-push
dev-push:
	docker push moov/ach-web-viewer:${DEV_VERSION}

# Extra utilities not needed for building

run: update build
	./bin/ach-web-viewer

docker-run:
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm moov-io/ach-web-viewer:${VERSION}

test: update
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

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@

dist: clean build
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=1 GOOS=windows go build -o bin/ach-web-viewer.exe cmd/ach-web-viewer/*
else
	CGO_ENABLED=1 GOOS=$(PLATFORM) go build -o bin/ach-web-viewer-$(PLATFORM)-amd64 cmd/ach-web-viewer/*
endif

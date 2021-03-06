# Makefile
BINARY       = $(shell basename "$(PWD)")
DOCKER_IMAGE = cube:latest
GIT_COMMIT   = github.com/shohi/cube/cmd/version.gitCommit
GIT_REVISION = $(shell git rev-parse --short HEAD)
GOENV        = CGO_ENABLED=0 GO111MODULE=on

default:
	@echo "$(BINARY) - $(GIT_REVISION)"

build:
	@$(GOENV) go build -ldflags "-X $(GIT_COMMIT)=$(GIT_REVISION)" -o "$(BINARY)"

install:
	@$(GOENV) go install -ldflags "-X $(GIT_COMMIT)=$(GIT_REVISION)"

docker:
	docker build \
		-t "${DOCKER_IMAGE}" \
		-f Dockerfile\
		.

.PHONY: default build install docker

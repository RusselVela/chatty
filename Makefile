# Default TAG when pushing Docker Image to dockerhub
TAG ?= latest

# Docker for mac and Linux needs specific arguments to mount ssh agent sock
ifeq ($(OS),Windows_NT)
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        DOCKER_EXTRA_ARGS=-v ${SSH_AUTH_SOCK}:${SSH_AUTH_SOCK} -e SSH_AUTH_SOCK=${SSH_AUTH_SOCK}
    endif
    ifeq ($(UNAME_S),Darwin)
        DOCKER_EXTRA_ARGS=-v /run/host-services/ssh-auth.sock:/run/host-services/ssh-auth.sock -e SSH_AUTH_SOCK=/run/host-services/ssh-auth.sock
    endif
endif

############ Build rules ############
.PHONY: init-vars
init-vars:
	$(eval MODULE_FQDN=$(shell go list -m))

.PHONY: generate
generate: init-vars
	go generate $(MODULE_FQDN)/...

.PHONY: build
build: generate
	go mod download
	go install $(MODULE_FQDN)/...

.PHONY: build-linux
build-linux: generate
	go mod download
	GOOS=linux GOARCH=amd64 go build -o output/bin/$(notdir $(CURDIR)) ./cmd/$(notdir $(CURDIR))

############ Image rules ############
.PHONY: buildx
buildx:
	docker buildx build ${BUILDX_OUTPUT} \
		--file build/Dockerfile \
		${BUILDX_EXTRA_ARGS} \
		--builder default .

.PHONY: image
image: buildx

.PHONY: push
push: BUILDX_OUTPUT=--output type=image,name=russelvela/chatty:${TAG},push=true --metadata-file=build.json
push: buildx

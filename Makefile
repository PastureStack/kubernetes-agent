TARGETS := $(shell ls scripts)

DAPPER_IMAGE ?= pasturestack-kubernetes-agent-dapper:go1.26.6
DAPPER_HOST_ARCH ?= amd64
DOCKER_VERSION ?= 29.6.2

.PHONY: $(TARGETS) deps trash trash-keep dapper-image

dapper-image:
	docker build \
		$(if $(DOCKER_BUILD_NETWORK),--network $(DOCKER_BUILD_NETWORK),) \
		--build-arg DAPPER_HOST_ARCH=$(DAPPER_HOST_ARCH) \
		--build-arg DOCKER_VERSION=$(DOCKER_VERSION) \
		-t $(DAPPER_IMAGE) \
		-f Dockerfile.dapper .

$(TARGETS): dapper-image
	docker run --rm \
		-v $(CURDIR):/source \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-e DAPPER_UID=$$(id -u) \
		-e DAPPER_GID=$$(id -g) \
		-e ARCH=$(DAPPER_HOST_ARCH) \
		-e IMAGE_NAME \
		-e TAG \
		-e VERSION_OVERRIDE \
		-e DOCKER_BUILD_NETWORK \
		-e GOFLAGS=-mod=mod \
		-e GOTOOLCHAIN=local \
		$(DAPPER_IMAGE) $@

trash:
	@echo "Legacy trash vendoring was removed; use Go modules."

trash-keep: trash

deps:
	go mod download

.DEFAULT_GOAL := ci

DOCKER_NAMESPACE = dnabic
DOCKER_REPOSITORY = ebscli
EBSCLI_VERSION ?= v0.0.5

lint:
	golint

build:
	@sh -c "$(CURDIR)/scripts/build.sh"

dev:
	@TF_DEV=1 sh -c "$(CURDIR)/scripts/build.sh"

style:
	gofmt -w .

test:
	go test

docker: build
	docker build -t $(DOCKER_NAMESPACE)/$(DOCKER_REPOSITORY):$(EBSCLI_VERSION) .

.PHONY: lint build dev style test docker

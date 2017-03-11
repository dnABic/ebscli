NAME = dnabic/ebscli
VERSION = v0.0.3

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
	docker build -t $(NAME):$(VERSION) .

.PHONY: lint build dev style test docker

USERNAME := grrrvahrrr
APP_NAME := gobelvl2
VERSION := latest
PROJECT := gobelvl2
GIT_COMMIT := $(shell git rev-parse HEAD)

all: run
run:
	go install -ldflags="-X '$(PROJECT)/version.Version=$(VERSION)' -X '$(PROJECT)/version.Commit=$(GIT_COMMIT)'"

build_container:
	docker build --build-arg=GIT_COMMIT=$(GIT_COMMIT) --build-arg=VERSION=$(VERSION) --build-arg=PROJECT=$(PROJECT) -t docker.io/$(USERNAME)/$(APP_NAME):$(VERSION) .

run_container:
	docker run -p 8080:8080 docker.io/$(USERNAME)/$(APP_NAME):$(VERSION)

push_container:
	docker push docker.io/$(USERNAME)/$(APP_NAME):$(VERSION)

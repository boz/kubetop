
DOCKER_IMAGE ?= kubetop
DOCKER_REPO  ?= abozanich/$(DOCKER_IMAGE)
DOCKER_TAG   ?= latest

IMG_LDFLAGS := -w -linkmode external -extldflags "-static"

build:
	go build

ifeq ($(shell uname -s),Darwin)
build-linux:
	GOOS=linux GOARCH=amd64 go build -o kubetop-linux
else
build-linux:
	CC=$$(which musl-gcc) go build --ldflags '$(IMG_LDFLAGS)' -o kubetop-linux
endif

test:
	govendor test +local

test-full:
	govendor test -v -race +local

image: build-linux
	docker build -t $(DOCKER_IMAGE) .

image-minikube: build-linux
	eval $$(minikube docker-env) && docker build -t $(DOCKER_IMAGE) .

image-push: image
	docker tag $(DOCKER_IMAGE) $(DOCKER_REPO):$(DOCKER_TAG)
	docker push $(DOCKER_REPO):$(DOCKER_TAG)

install-libs:
	govendor install +vendor,^program

install-deps:
	go get github.com/kardianos/govendor
	govendor sync

clean:
	rm kubetop kubetop-linux 2>/dev/null || true

.PHONY: build build-linux \
	test test-full \
	image image-minikube image-push \
	install-libs install-deps \
	clean

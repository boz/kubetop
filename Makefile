IMG_LDFLAGS := -w -linkmode external -extldflags "-static"

DOCKER_IMAGE ?= kubetop
DOCKER_REPO  ?= abozanich/$(DOCKER_IMAGE)
DOCKER_TAG   ?= latest

build:
	go build

build-linux:
	CC=$$(which musl-gcc) go build --ldflags '$(IMG_LDFLAGS)' -o kubetop-linux

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

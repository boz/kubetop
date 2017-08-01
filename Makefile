IMG_LDFLAGS := -w -linkmode external -extldflags "-static"

build:
	go build

build-linux:
	CC=$$(which musl-gcc) go build --ldflags '$(IMG_LDFLAGS)' -o kubetop-linux

test:
	govendor test +local

test-full:
	govendor test -v -race +local

image: build-linux
	docker build -t kubetop .

image-minikube: build-linux
	eval $$(minikube docker-env) && docker build -t kubetop .

install-libs:
	govendor install +vendor,^program

install-deps:
	go get github.com/kardianos/govendor
	govendor sync

clean:
	rm kubetop kubetop-linux 2>/dev/null || true

.PHONY: build build-linux \
	test test-full \
	image image-minikube \
	install-libs install-deps \
	clean

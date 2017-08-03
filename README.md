# kubetop: kubernetes cli console [![Build Status](https://travis-ci.org/boz/kubetop.svg?&branch=master)](https://travis-ci.org/boz/kubetop)

## Building

Install build & dev dependencies:

  * [govendor](https://github.com/kardianos/govendor)
  * [minikube](https://kubernetes.io/docs/getting-started-guides/minikube/)
  * _linux only_: [musl-gcc](https://www.musl-libc.org/how.html) for building docker images.

Install source code and golang dependencies:

```sh
go get -d github.com/boz/kubetop
cd $GOPATH/src/github.com/boz/kubetop
make install-deps
```

build `kubetop` binary:

```sh
make
```

build `kubetop` docker image on `minikube`'s docker

```sh
make image-minikube
```

## Running

Use `kubectl`'s default context:

```sh
./kubetop
```

Run as a kubernetes job:

```sh
kubectl run -it --rm --restart=Never --image=abozanich/kubetop kubetop
```

Build an image locally and test:

```sh
make image-minikube
kubectl run -it --rm --restart=Never --image=kubetop kubetop
```

There are kube definitions in [_example](https://github.com/boz/kubetop/tree/master/_example) for dev testing

```sh
ls _example/*.yml | xargs -n 1 kubectl create -f
```

## Status

Work in progress

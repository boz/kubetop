# kubetop: kubernetes cli console [![Build Status](https://travis-ci.com/boz/kubetop.svg?token=xMx9pPujMteGc5JpGjzX&branch=master)](https://travis-ci.com/boz/kubetop)

## Building

```sh
go get -d github.com/boz/kubetop
cd $GOPATH/src/github.com/boz/kubetop
make install-deps
make
```

## Running

```sh
./kubetop
```

```sh
make image-minikube
kubectl run -it --rm --restart=Never --image=kubetop kubetop
```

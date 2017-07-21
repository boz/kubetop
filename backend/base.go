package backend

import (
	"github.com/boz/kcache/types/node"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
)

type BasePodController interface {
	pod.CacheController
	pod.Publisher
	Done() <-chan struct{}
}

type BaseServiceController interface {
	service.CacheController
	service.Publisher
	Done() <-chan struct{}
}

type BaseNodeController interface {
	node.CacheController
	node.Publisher
	Done() <-chan struct{}
}

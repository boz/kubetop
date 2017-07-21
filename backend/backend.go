package backend

import (
	"context"
	"sync"

	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/kubernetes"
)

type Backend interface {
	Pods() (BasePodController, error)
	Services() (BaseServiceController, error)
	Close()
}

type backend struct {
	clientset kubernetes.Interface

	pods     pod.Controller
	services service.Controller

	ctx    context.Context
	cancel context.CancelFunc
	env    util.Env
}

func NewBackend(env util.Env, clientset kubernetes.Interface) Backend {
	ctx, cancel := context.WithCancel(context.TODO())
	return &backend{
		clientset: clientset,
		ctx:       ctx,
		cancel:    cancel,
		env:       env.ForComponent("backend/backend"),
	}
}

func (b *backend) Close() {
	var wg sync.WaitGroup
	b.cancel()

	b.env.Log().Debug("stopping...")

	b.doClose(&wg, b.pods)
	b.doClose(&wg, b.services)

	wg.Wait()
}

func (b *backend) Pods() (BasePodController, error) {
	if b.pods == nil {
		controller, err := pod.NewController(b.ctx, b.env.Logutil(), b.clientset, "")
		if err != nil {
			return nil, err
		}
		b.pods = controller
	}
	return b.pods, nil
}

func (b *backend) Services() (BaseServiceController, error) {
	if b.services == nil {
		controller, err := service.NewController(b.ctx, b.env.Logutil(), b.clientset, "")
		if err != nil {
			return nil, err
		}
		b.services = controller
	}
	return b.services, nil
}

type closeable interface {
	Done() <-chan struct{}
}

func (b *backend) doClose(wg *sync.WaitGroup, db closeable) {
	if db == nil {
		return
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-db.Done()
	}()
}

package backend

import (
	"context"
	"sync"

	lr "github.com/boz/go-logutil/logrus"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/kubernetes"
)

type Backend interface {
	Pods() (pod.Controller, error)
	Services() (service.Controller, error)

	Stop()
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

func (b *backend) Stop() {
	var wg sync.WaitGroup
	b.cancel()

	b.env.Log().Debug("stopping...")

	b.doStop(&wg, b.pods)
	b.doStop(&wg, b.services)

	wg.Wait()
}

func (b *backend) Pods() (pod.Controller, error) {
	if b.pods == nil {
		log := lr.New(b.env.Log())
		controller, err := pod.NewController(b.ctx, log, b.clientset, "")
		if err != nil {
			return nil, err
		}
		b.pods = controller
	}
	return b.pods, nil
}

func (b *backend) Services() (service.Controller, error) {
	if b.services == nil {
		log := lr.New(b.env.Log())
		controller, err := service.NewController(b.ctx, log, b.clientset, "")
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

func (b *backend) doStop(wg *sync.WaitGroup, db closeable) {
	if db == nil {
		return
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-db.Done()
	}()
}

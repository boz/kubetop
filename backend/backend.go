package backend

import (
	"sync"

	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/kubernetes"
)

type Backend interface {
	Pods() (pod.BaseDatasource, error)

	Stop()
}

type backend struct {
	clientset kubernetes.Interface

	pods pod.Datasource

	env util.Env
}

func NewBackend(env util.Env, clientset kubernetes.Interface) Backend {
	return &backend{
		clientset: clientset,
		env:       env.ForComponent("backend/backend"),
	}
}

func (b *backend) Stop() {
	var wg sync.WaitGroup

	b.env.Log().Debug("stopping...")

	b.doStop(&wg, b.pods)

	wg.Wait()
}

func (b *backend) Pods() (pod.BaseDatasource, error) {
	if b.pods == nil {
		pods, err := pod.NewBase(b.env, b.clientset)
		if err != nil {
			return nil, err
		}
		b.pods = pods
	}
	return b.pods, nil
}

type closeable interface {
	Close()
}

func (b *backend) doStop(wg *sync.WaitGroup, db closeable) {
	if db == nil {
		return
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		db.Close()
	}()
}

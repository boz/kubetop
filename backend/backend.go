package backend

import (
	"sync"

	"github.com/boz/kubetop/backend/pod"
	"k8s.io/client-go/kubernetes"
)

type Backend interface {
	Pods(pod.Filters) (pod.Datasource, error)

	Stop()
}

type backend struct {
	clientset kubernetes.Interface

	pods pod.Database
}

func NewBackend(clientset kubernetes.Interface) Backend {
	return &backend{
		clientset: clientset,
	}
}

func (b *backend) Stop() {
	var wg sync.WaitGroup
	b.doStop(&wg, b.pods)

	wg.Wait()
}

func (b *backend) Pods(filters pod.Filters) (pod.Datasource, error) {
	if b.pods == nil {
		pods, err := pod.NewDatabase(b.clientset)
		if err != nil {
			return nil, err
		}
		b.pods = pods
	}
	return b.pods.Filter(filters), nil
}

type stopper interface {
	Stop()
}

func (b *backend) doStop(wg *sync.WaitGroup, db stopper) {
	if db == nil {
		return
	}
	wg.Add(1)

	go func() {
		defer wg.Done()
		db.Stop()
	}()
}

package pod

import (
	"context"

	lr "github.com/boz/go-logutil/logrus"
	"github.com/boz/kcache"
	"github.com/boz/kcache/client"
	"github.com/boz/kubetop/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

type BaseDatasource interface {
	Get(Pod) (Pod, error)
	List() ([]Pod, error)
	Subscribe() Subscription
	Filter(Filters) Datasource
}

type Datasource interface {
	BaseDatasource
	Close()
}

type Subscription interface {
	Get(Pod) (Pod, error)
	List() ([]Pod, error)
	Events() <-chan Event
	Close()
	Closed() <-chan struct{}
}

type Filter interface {
	Accept(Pod) bool
}

type Filters []Filter

func NewBase(env util.Env, clientset kubernetes.Interface) (Datasource, error) {
	env = env.ForComponent("backend/pod/base-datasource")
	env = env.WithFields("model", "pod")

	client := client.ForResource(clientset.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, fields.Everything())
	ctx := context.Background()
	log := lr.New(env.Log().WithField("layer", "controller"))

	controller, err := kcache.NewController(ctx, log, client)
	if err != nil {
		return nil, err
	}
	return &_datasource{controller, newAdapter(env), env}, nil
}

type _datasource struct {
	controller kcache.Controller
	adapter    *adapter
	env        util.Env
}

func (db *_datasource) Filter(filters Filters) Datasource {
	env := db.env.ForComponent("backend/pod/datasource").WithID()
	log := lr.New(env.Log())
	controller := kcache.CloneWithFilter(log, db.controller, kcache.NullFilter())
	return &_datasource{controller, db.adapter, env}
}

func (db *_datasource) Stop() {
	db.controller.Close()
}

func (ds *_datasource) Get(p Pod) (Pod, error) {
	return doGet(ds.adapter, ds.controller.Cache(), p)
}

func (ds *_datasource) List() ([]Pod, error) {
	return doList(ds.adapter, ds.controller.Cache())
}

func (ds *_datasource) Subscribe() Subscription {
	parent := ds.controller.Subscribe()
	return newSubscription(ds.env, ds, parent)
}

func (ds *_datasource) Close() {
	ds.controller.Close()
}

type subscription struct {
	env     util.Env
	adapter Adapter
	parent  kcache.Subscription
	outch   chan Event
}

func newSubscription(env util.Env, ds Datasource, parent kcache.Subscription) *subscription {
	env = env.ForComponent("backend/pod/subcription").WithID()
	s := &subscription{
		env:     env,
		adapter: newAdapter(env),
		parent:  parent,
		outch:   make(chan Event),
	}
	go s.translateEvents()
	return s
}

func (s *subscription) Close() {
	s.parent.Close()
}

func (s *subscription) Closed() <-chan struct{} {
	return s.parent.Done()
}

func (s *subscription) Get(p Pod) (Pod, error) {
	return doGet(s.adapter, s.parent.Cache(), p)
}

func (s *subscription) List() ([]Pod, error) {
	return doList(s.adapter, s.parent.Cache())
}

func (s *subscription) Events() <-chan Event {
	return s.outch
}

func (s *subscription) translateEvents() {
	defer close(s.outch)
	for ev := range s.parent.Events() {
		pod, err := s.adapter.FromResource(ev.Resource())
		if err != nil {
			s.env.LogErr(err, "adapt event")
			continue
		}
		s.outch <- newEvent(ev.Type(), pod)
	}
}

func doGet(adapter Adapter, cr kcache.CacheReader, p Pod) (Pod, error) {
	obj, err := cr.GetObject(p.Resource())
	if err != nil {
		return nil, err
	}
	return adapter.FromResource(obj)
}

func doList(adapter Adapter, cr kcache.CacheReader) ([]Pod, error) {
	objs, err := cr.List()
	if err != nil {
		return nil, err
	}
	return adapter.FromResourceList(objs)
}

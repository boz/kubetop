package service

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
	Ready() <-chan struct{}
	List() ([]Service, error)
	Subscribe(kcache.Filter) Subscription
	Filter(kcache.Filter) Datasource
}

type Datasource interface {
	BaseDatasource
	Refilter(kcache.Filter)
	Close()
}

type Subscription interface {
	Ready() <-chan struct{}
	List() ([]Service, error)
	Events() <-chan Event
	Refilter(kcache.Filter)
	Close()
	Closed() <-chan struct{}
}

func NewBase(env util.Env, clientset kubernetes.Interface) (Datasource, error) {
	env = env.ForComponent("backend/service/base-datasource")
	env = env.WithFields("model", "service")

	client := client.ForResource(clientset.CoreV1().RESTClient(), "services", metav1.NamespaceAll, fields.Everything())
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

type _filterDatasource struct {
	parent kcache.FilterController
	_datasource
}

func (db *_datasource) Filter(filter kcache.Filter) Datasource {
	env := db.env.ForComponent("backend/service/datasource").WithID()
	log := lr.New(env.Log())
	controller := kcache.CloneWithFilter(log, db.controller, filter)
	return &_filterDatasource{controller, _datasource{controller, db.adapter, env}}
}

func (db *_datasource) Ready() <-chan struct{} {
	return db.controller.Ready()
}

func (db *_datasource) Stop() {
	db.controller.Close()
}

func (ds *_datasource) List() ([]Service, error) {
	return doList(ds.adapter, ds.controller.Cache())
}

func (ds *_datasource) Subscribe(filter kcache.Filter) Subscription {
	return newSubscription(ds.env, ds, filter)
}

func (ds *_datasource) Close() {
	ds.controller.Close()
}

func (ds *_datasource) Refilter(filter kcache.Filter) {
	panic("not implemented")
}

func (ds *_filterDatasource) Refilter(filter kcache.Filter) {
	ds.parent.Refilter(filter)
}

type subscription struct {
	env     util.Env
	adapter Adapter
	parent  kcache.FilterSubscription
	outch   chan Event
}

func newSubscription(env util.Env, ds *_datasource, filter kcache.Filter) *subscription {
	env = env.ForComponent("backend/service/subcription").WithID()
	log := lr.New(ds.env.Log())
	parent := kcache.SubscribeWithFilter(log, ds.controller, filter)
	s := &subscription{
		env:     env,
		adapter: newAdapter(env),
		parent:  parent,
		outch:   make(chan Event),
	}
	go s.translateEvents()
	return s
}

func (s *subscription) Ready() <-chan struct{} {
	return s.parent.Ready()
}

func (s *subscription) Close() {
	s.parent.Close()
}

func (s *subscription) Closed() <-chan struct{} {
	return s.parent.Done()
}

func (s *subscription) List() ([]Service, error) {
	return doList(s.adapter, s.parent.Cache())
}

func (s *subscription) Events() <-chan Event {
	return s.outch
}

func (s *subscription) Refilter(filter kcache.Filter) {
	s.parent.Refilter(filter)
}

func (s *subscription) translateEvents() {
	defer close(s.outch)
	for ev := range s.parent.Events() {
		service, err := s.adapter.FromResource(ev.Resource())
		if err != nil {
			s.env.LogErr(err, "adapt event")
			continue
		}
		s.outch <- newEvent(ev.Type(), service)
	}
}

func doList(adapter Adapter, cr kcache.CacheReader) ([]Service, error) {
	objs, err := cr.List()
	if err != nil {
		return nil, err
	}
	return adapter.FromResourceList(objs)
}

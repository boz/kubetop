package pod

import (
	"github.com/boz/kubetop/backend/database"
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type Database interface {
	Filter(Filters) Datasource
	Stop()
}

type Datasource interface {
	Get(Pod) (Pod, error)
	List() ([]Pod, error)
	Subscribe() Subscription
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

func NewDatabase(env util.Env, clientset kubernetes.Interface) (Database, error) {
	env = env.ForComponent("backend/pod/database")
	env = env.WithFields("model", "pod")

	lw := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(), "pods", api.NamespaceAll, fields.Everything())

	db, err := database.NewDatabase(
		env, lw, &v1.Pod{}, database.DefaultResyncPeriod, database.BaseIndexers())

	if err != nil {
		return nil, err
	}
	return &_database{db, newAdapter(env), env}, nil
}

type _database struct {
	db database.Database

	adapter *adapter

	env util.Env
}

func (db *_database) Filter(filters Filters) Datasource {
	return newDatasource(db.env, db.adapter, db.db)
}

func (db *_database) Stop() {
	db.db.Stop()
}

type datasource struct {
	db      database.Database
	adapter *adapter

	env util.Env
}

func newDatasource(env util.Env, adapter *adapter, db database.Database) *datasource {
	env = env.ForComponent("backend/pod/datasource")
	env = env.WithID()
	return &datasource{db, adapter, env}
}

func (ds *datasource) Get(p Pod) (Pod, error) {
	obj, exists, err := ds.db.Indexer().Get(p.Resource())
	switch {
	case err != nil:
		return nil, err
	case !exists:
		return nil, database.ErrNotFound
	default:
		return ds.adapter.FromResource(obj)
	}
}

func (ds *datasource) List() ([]Pod, error) {
	objs := ds.db.Indexer().List()
	return ds.adapter.FromResourceList(objs)
}

func (ds *datasource) Subscribe() Subscription {
	parent := ds.db.Subscribe()
	return newSubscription(ds.env, ds, parent)
}

type subscription struct {
	env     util.Env
	ds      Datasource
	adapter Adapter
	parent  database.Subscription
	outch   chan Event
}

func newSubscription(env util.Env, ds Datasource, parent database.Subscription) *subscription {
	env = env.ForComponent("backend/pod/subcription").WithID()
	s := &subscription{
		env:     env,
		ds:      ds,
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
	return s.parent.Closed()
}

func (s *subscription) Get(p Pod) (Pod, error) {
	return s.ds.Get(p)
}

func (s *subscription) List() ([]Pod, error) {
	return s.ds.List()
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

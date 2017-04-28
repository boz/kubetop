package pod

import (
	"fmt"

	"github.com/boz/kubetop/backend/database"
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type Adapter interface {
	FromResource(interface{}) (Pod, error)
	ToResource(Pod) (*v1.Pod, error)
}

type Database interface {
	Filter(Filters) Datasource
	Stop()
}

type Pod interface {
	Resource() *v1.Pod
}

type Event interface {
	Resource() Pod
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

	pods := make([]Pod, 0, len(objs))

	for _, obj := range objs {
		pod, err := ds.adapter.FromResource(obj)
		if err != nil {
			return nil, err
		}
		pods = append(pods, pod)
	}

	return pods, nil
}

func (ds *datasource) Subscribe() Subscription {
	return nil
}

type adapter struct {
	env util.Env
}

func newAdapter(env util.Env) *adapter {
	return &adapter{env}
}

func (a adapter) ToResource(p Pod) (*v1.Pod, error) {
	return p.Resource(), nil
}

func (a adapter) FromResource(obj interface{}) (Pod, error) {
	switch obj := obj.(type) {
	case *v1.Pod:
		return newPod(a.env, obj), nil
	default:
		return nil, fmt.Errorf("invalid type: %T", obj)
	}
}

type pod struct {
	resource *v1.Pod
	env      util.Env
}

func newPod(env util.Env, resource *v1.Pod) *pod {
	return &pod{resource, env}
}

func (p *pod) Resource() *v1.Pod {
	return p.resource
}

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

type Pod interface {
}

type Event interface {
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
	return &_database{db, env}, nil
}

type _database struct {
	db database.Database

	env util.Env
}

func (db *_database) Filter(filters Filters) Datasource {
	return nil
}

func (db *_database) Stop() {
	db.db.Stop()
}

package database

import (
	"time"

	"github.com/boz/kubetop/util"

	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

const (
	DefaultResyncPeriod = time.Second
)

func BaseIndexers() cache.Indexers {
	return map[string]cache.IndexFunc{
		cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
	}
}

type Database interface {
	Subscribe() Subscription
	Stop()
}

type Event interface {
	Resource() interface{}
}

type EventUpdate struct{ event }
type EventCreate struct{ event }
type EventDelete struct{ event }

type event struct {
	resource interface{}
}

func (e event) Resource() interface{} {
	return e.resource
}

type database struct {
	controller *cache.Controller
	indexer    cache.Indexer

	subscribers map[*subscription]struct{}

	events chan Event

	subch   chan chan<- Subscription
	unsubch chan *subscription

	stopch  chan struct{}
	cdonech chan struct{}
	donech  chan struct{}

	env util.Env
}

func NewDatabase(
	env util.Env,
	lw cache.ListerWatcher,
	obj runtime.Object,
	period time.Duration,
	indexers cache.Indexers,
) (Database, error) {

	env = env.ForComponent("backend/database/database")

	db := &database{
		subscribers: make(map[*subscription]struct{}),

		events: make(chan Event),

		// create new subscription
		subch: make(chan chan<- Subscription),

		// unsubscribe
		unsubch: make(chan *subscription),

		// start shutdown
		stopch: make(chan struct{}),

		// controller done
		cdonech: make(chan struct{}),

		// completely shut down
		donech: make(chan struct{}),

		env: env,
	}

	handlers := cache.ResourceEventHandlerFuncs{
		db.onResourceAdd,
		db.onResourceUpdate,
		db.onResourceDelete,
	}

	db.indexer, db.controller = cache.NewIndexerInformer(lw, obj, period, handlers, indexers)

	go func() {
		db.controller.Run(db.stopch)
		db.cdonech <- struct{}{}
	}()

	return db, nil
}

func (db *database) run() {
	defer close(db.donech)

	stopch := db.stopch
	subch := db.subch

	for {
		select {
		case <-stopch:
			stopch = nil
			subch = nil
			db.subscribers = make(map[*subscription]struct{})
		case <-db.cdonech:
			return

		case ch := <-subch:

			// make new subscription and send it back

			if sub := db.subscribe(ch); sub != nil {
				db.subscribers[sub] = struct{}{}
			}

		case ev := <-db.events:
			for s, _ := range db.subscribers {
				s.postEvent(ev)
			}
		}
	}
}

func (db *database) onResourceAdd(obj interface{}) {
	select {
	case <-db.donech:
	case db.events <- EventCreate{event{obj}}:
	}
}

func (db *database) onResourceUpdate(prev, cur interface{}) {
	select {
	case <-db.donech:
	case db.events <- EventUpdate{event{cur}}:
	}
}

func (db *database) onResourceDelete(obj interface{}) {
	select {
	case <-db.donech:
	case db.events <- EventDelete{event{obj}}:
	}
}

func (db *database) subscribe(ch chan<- Subscription) *subscription {
	s := newSubscriptionForDB(db.env, db)
	select {
	case <-db.donech:
		return nil
	case ch <- s:
		return s
	}
}

func (db *database) unsubscribe(s *subscription) {
	select {
	case <-db.donech:
	case db.unsubch <- s:
	}
}

func (db *database) Subscribe() Subscription {
	ch := make(chan Subscription)
	subch := db.subch
	for {
		select {
		case <-db.donech:
			return stoppedSubscription(db)
		case subch <- ch:
			subch = nil
		case sub := <-ch:
			if sub == nil {
				return stoppedSubscription(db)
			}
			return sub
		}
	}
}

func (db *database) Stop() {
	for {
		select {
		case <-db.donech:
			return
		case db.stopch <- struct{}{}:
		}
	}
}

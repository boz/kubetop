package database

import (
	"fmt"
	"time"

	"github.com/boz/kubetop/util"

	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

const (
	DefaultResyncPeriod     = time.Second
	syncPerformedPollPeriod = time.Second / 10
)

var (
	ErrNotFound = fmt.Errorf("Not found")
)

func BaseIndexers() cache.Indexers {
	return map[string]cache.IndexFunc{
		cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
	}
}

type Database interface {
	Indexer() cache.Indexer
	Subscribe() Subscription
	Stop()

	Ready() <-chan struct{}
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

	// incoming events
	events chan Event

	// create new subscription
	subch chan chan<- Subscription

	// unsubscribe
	unsubch chan *subscription

	// closed when controller has synced
	readych chan struct{}

	// initiate stop
	stopch chan struct{}

	// closed when stopping
	stoppingch chan struct{}

	// closed when controller done
	cdonech chan struct{}

	// closed when completely shutdown
	donech chan struct{}

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
		events:      make(chan Event),
		subch:       make(chan chan<- Subscription),
		unsubch:     make(chan *subscription),
		readych:     make(chan struct{}),
		stopch:      make(chan struct{}),
		stoppingch:  make(chan struct{}),
		cdonech:     make(chan struct{}),
		donech:      make(chan struct{}),
		env:         env,
	}

	handlers := cache.ResourceEventHandlerFuncs{
		db.onResourceAdd,
		db.onResourceUpdate,
		db.onResourceDelete,
	}

	db.indexer, db.controller = cache.NewIndexerInformer(lw, obj, period, handlers, indexers)

	go db.run()
	go db.pollReady()
	go db.runController()

	return db, nil
}

// close readych when controller has synced.
func (db *database) pollReady() {
	for {
		if db.controller.HasSynced() {
			close(db.readych)
			return
		}
		select {
		case <-db.stoppingch:
			// shut down before ready
			return
		case <-time.After(syncPerformedPollPeriod):
			// retry
		}
	}
}

// run controler. write to cdonech on controller exit
// XXX: controller.Run() never returns
func (db *database) runController() {
	db.controller.Run(db.stoppingch)
	db.cdonech <- struct{}{}
}

func (db *database) run() {
	defer close(db.donech)

	stopch := db.stopch
	subch := db.subch

	for {
		select {
		case <-stopch:
			// begin shutdown

			close(db.stoppingch)

			stopch = nil
			subch = nil

			db.subscribers = make(map[*subscription]struct{})

			// XXX: return here as cdonech is never closed.
			return

		case <-db.cdonech:
			// controller is done; exit

			return

		case ch := <-subch:
			// make new subscription and send it back

			if sub := db.subscribe(ch); sub != nil {
				db.subscribers[sub] = struct{}{}
			}

		case ev := <-db.events:
			// forward incoming events

			for s, _ := range db.subscribers {
				s.postEvent(ev)
			}
		}
	}
}

func (db *database) onResourceAdd(obj interface{}) {
	select {
	case <-db.stoppingch:
	case db.events <- EventCreate{event{obj}}:
	}
}

func (db *database) onResourceUpdate(prev, cur interface{}) {
	select {
	case <-db.stoppingch:
	case db.events <- EventUpdate{event{cur}}:
	}
}

func (db *database) onResourceDelete(obj interface{}) {
	select {
	case <-db.stoppingch:
	case db.events <- EventDelete{event{obj}}:
	}
}

func (db *database) subscribe(ch chan<- Subscription) *subscription {
	s := newSubscriptionForDB(db.env, db)
	select {
	case <-db.stoppingch:
		close(ch)
		return nil
	case ch <- s:
		return s
	}
}

func (db *database) unsubscribe(s *subscription) {
	select {
	case <-db.stoppingch:
	case db.unsubch <- s:
	}
}

func (db *database) Subscribe() Subscription {
	ch := make(chan Subscription)
	subch := db.subch
	for {
		select {
		case <-db.stoppingch:
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

func (db *database) Ready() <-chan struct{} {
	return db.readych
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

func (db *database) Indexer() cache.Indexer {
	return db.indexer
}

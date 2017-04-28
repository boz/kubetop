package database

import (
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/tools/cache"
)

type Subscription interface {
	Indexer() cache.Indexer

	Events() <-chan Event

	Close()
	Closed() <-chan struct{}
}

type subscription struct {
	db *database

	buffer []Event

	inch  chan Event
	outch chan Event

	stopch     chan struct{}
	stoppingch chan struct{}
	donech     chan struct{}

	env util.Env
}

func newSubscriptionForDB(env util.Env, db *database) *subscription {
	s := &subscription{
		db:         db,
		inch:       make(chan Event),
		outch:      make(chan Event),
		stopch:     make(chan struct{}),
		stoppingch: make(chan struct{}),
		donech:     make(chan struct{}),
		env:        env.WithID(),
	}
	go s.run()
	return s
}

func (s *subscription) run() {
	defer close(s.donech)

	for {
		var outch chan Event
		var head Event

		if len(s.buffer) > 0 {
			head = s.buffer[0]
			outch = s.outch
		}

		select {
		case <-s.db.stoppingch:
			// db is stopping; exit without unsubscribe
			close(s.stoppingch)
			return
		case <-s.stopch:
			close(s.stoppingch)
			s.db.unsubscribe(s)
			return
		case ev := <-s.inch:
			s.buffer = append(s.buffer, ev)
		case outch <- head:
			s.buffer = s.buffer[1:]
		}
	}
}

func (s *subscription) postEvent(ev Event) {
	select {
	case <-s.stoppingch:
	case s.inch <- ev:
	}
}

func (s *subscription) Indexer() cache.Indexer {
	return s.db.indexer
}

func (s *subscription) Events() <-chan Event {
	return s.outch
}

func (s *subscription) Close() {
	for {
		select {
		case <-s.stoppingch:
			<-s.donech
			return
		case s.stopch <- struct{}{}:
		}
	}
}

func (s *subscription) Closed() <-chan struct{} {
	return s.donech
}

func stoppedSubscription(db *database) Subscription {
	s := &subscription{
		db:         db,
		stoppingch: make(chan struct{}),
		donech:     make(chan struct{}),
	}
	close(s.stoppingch)
	close(s.donech)
	return s
}
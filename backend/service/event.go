package service

import "github.com/boz/kcache"

type Event interface {
	Type() kcache.EventType
	Resource() Service
}

type event struct {
	etype    kcache.EventType
	resource Service
}

func (ev *event) Type() kcache.EventType { return ev.etype }
func (ev *event) Resource() Service      { return ev.resource }

func newEvent(et kcache.EventType, r Service) Event {
	return &event{et, r}
}

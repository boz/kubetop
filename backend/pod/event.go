package pod

import "github.com/boz/kcache"

type Event interface {
	Type() kcache.EventType
	Resource() Pod
}

type event struct {
	etype    kcache.EventType
	resource Pod
}

func (ev *event) Type() kcache.EventType { return ev.etype }
func (ev *event) Resource() Pod          { return ev.resource }

func newEvent(et kcache.EventType, r Pod) Event {
	return &event{et, r}
}

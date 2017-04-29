package pod

import "github.com/boz/kubetop/backend/database"

type Event interface {
	Type() database.EventType
	Resource() Pod
}

type event struct {
	etype    database.EventType
	resource Pod
}

func (ev *event) Type() database.EventType { return ev.etype }
func (ev *event) Resource() Pod            { return ev.resource }

func newEvent(et database.EventType, r Pod) Event {
	return &event{et, r}
}

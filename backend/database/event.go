package database

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

type Event interface {
	Type() EventType
	Resource() interface{}
}

type EventUpdate struct{ event }
type EventCreate struct{ event }
type EventDelete struct{ event }

type event struct {
	etype    EventType
	resource interface{}
}

func (e event) Resource() interface{} {
	return e.resource
}

func (e event) Type() EventType {
	return e.etype
}

func NewEventCreate(resource interface{}) EventCreate {
	return EventCreate{event{EventTypeCreate, resource}}
}

func NewEventUpdate(resource interface{}) EventUpdate {
	return EventUpdate{event{EventTypeUpdate, resource}}
}

func NewEventDelete(resource interface{}) EventDelete {
	return EventDelete{event{EventTypeDelete, resource}}
}

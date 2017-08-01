package requests

import (
	"github.com/boz/kcache/nsname"
	"github.com/boz/kubetop/ui/elements"
)

const (
	PodIndexRoute     elements.Route = "/pod"
	PodShowRoute      elements.Route = "/pod/show"
	ServiceIndexRoute elements.Route = "/service"
	EventIndexRoute   elements.Route = "/event"
	NodeIndexRoute    elements.Route = "/node"
)

func PodIndexRequest() elements.Request {
	return elements.NewRequest(PodIndexRoute)
}

func PodShowRequest(id nsname.NSName) elements.NSNameRequest {
	return elements.NewNSNameRequest(PodShowRoute, id)
}

func ServiceIndexRequest() elements.Request {
	return elements.NewRequest(ServiceIndexRoute)
}

func EventIndexRequest() elements.Request {
	return elements.NewRequest(EventIndexRoute)
}

func NodeIndexRequest() elements.Request {
	return elements.NewRequest(NodeIndexRoute)
}

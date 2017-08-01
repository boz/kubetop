package requests

import "github.com/boz/kubetop/ui/elements"

const (
	PodIndexRoute     elements.Route = "/pod"
	ServiceIndexRoute elements.Route = "/service"
	EventIndexRoute   elements.Route = "/event"
	NodeIndexRoute    elements.Route = "/node"
)

func PodIndexRequest() elements.Request {
	return elements.NewRequest(PodIndexRoute)
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

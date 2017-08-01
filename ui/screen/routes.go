package screen

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/screen/event"
	"github.com/boz/kubetop/ui/screen/node"
	"github.com/boz/kubetop/ui/screen/pod"
	"github.com/boz/kubetop/ui/screen/requests"
	"github.com/boz/kubetop/ui/screen/service"
)

func RegisterRoutes(router elements.Router) {
	router.Register(requests.PodIndexRoute, elements.NewHandler(pod.NewIndex))
	router.Register(requests.ServiceIndexRoute, elements.NewHandler(service.NewIndex))
	router.Register(requests.NodeIndexRoute, elements.NewHandler(node.NewIndex))
	router.Register(requests.EventIndexRoute, elements.NewHandler(event.NewIndex))
}

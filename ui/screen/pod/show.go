package pod

import (
	"fmt"

	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	eview "github.com/boz/kubetop/ui/view/event"
	pview "github.com/boz/kubetop/ui/view/pod"
	sview "github.com/boz/kubetop/ui/view/service"
	"github.com/gdamore/tcell/views"
)

type showScreen struct {
	layout elements.Panes
	ctx    elements.Context
	views.WidgetWatchers
}

func NewShow(ctx elements.Context, req elements.NSNameRequest) (elements.Screen, error) {
	ctx = ctx.New("pod/show")

	ds, err := newShowDS(ctx, req.NSName())
	if err != nil {
		ctx.Env().LogErr(err, "opening data source")
		ctx.Close()
		return nil, err
	}

	layout := elements.NewVSections(ctx.Env(), true)

	containerw := showContainersSection(ctx, ds)
	layout.Append(containerw)

	servicesw := showServicesSection(ctx, ds)
	layout.Append(servicesw)

	eventsw := showEventsSection(ctx, ds)
	layout.Append(eventsw)

	return elements.NewScreen(ctx, req, fmt.Sprintf("Pod %v", req.NSName()), layout), nil
}

func showContainersSection(ctx elements.Context, ds *showDS) elements.Section {
	ctx = ctx.New("pod/show/containers")
	content := table.NewWidget(ctx.Env(), pview.ContainersTableColumns(), true)
	pview.MonitorUnitary(ctx, ds.pods, pview.NewContainersTable(content))
	return elements.NewSection(ctx.Env(), "Containers", content)
}

func showServicesSection(ctx elements.Context, ds *showDS) elements.Section {
	ctx = ctx.New("pod/show/services")
	content := table.NewWidget(ctx.Env(), sview.TableColumns(), true)
	sview.Monitor(ctx, ds.services, sview.NewTable(content))
	return elements.NewSection(ctx.Env(), "Services", content)
}

func showEventsSection(ctx elements.Context, ds *showDS) elements.Section {
	ctx = ctx.New("pod/show/events")
	content := table.NewWidget(ctx.Env(), eview.TableColumns(), true)
	eview.Monitor(ctx, ds.events, eview.NewTable(content))
	return elements.NewSection(ctx.Env(), "Events", content)
}

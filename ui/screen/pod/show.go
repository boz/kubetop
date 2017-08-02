package pod

import (
	"fmt"

	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
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

	layout := elements.NewVPanes(true)

	containerw := showContainersPanel(ctx, ds)
	layout.PushBackWidget(containerw)

	servicesw := showServicesPanel(ctx, ds)
	layout.PushBackWidget(servicesw)

	eventsw := showEventsPanel(ctx, ds)
	layout.PushBackWidget(eventsw)

	return elements.NewScreen(ctx, req, fmt.Sprintf("Pod %v", req.NSName()), layout), nil
}

func showContainersPanel(ctx elements.Context, ds *showDS) views.Widget {
	ctx = ctx.New("pod/show/containers")
	content := table.NewWidget(ctx.Env(), pview.ContainersTableColumns(), true)
	pview.MonitorUnitary(ctx, ds.pods, pview.NewContainersTable(content))
	return newPanel("Containers", content)
}

func showServicesPanel(ctx elements.Context, ds *showDS) views.Widget {
	ctx = ctx.New("pod/show/services")
	content := table.NewWidget(ctx.Env(), sview.TableColumns(), true)
	sview.Monitor(ctx, ds.services, sview.NewTable(content))
	return newPanel("Services", content)
}

func showEventsPanel(ctx elements.Context, ds *showDS) views.Widget {
	ctx = ctx.New("pod/show/events")
	content := table.NewWidget(ctx.Env(), eview.TableColumns(), true)
	eview.Monitor(ctx, ds.events, eview.NewTable(content))
	return newPanel("Events", content)
}

func newPanel(title string, content views.Widget) views.Widget {
	titlew := views.NewTextBar()
	titlew.SetCenter(title, theme.Base)
	layout := elements.NewVPanes(false)
	layout.PushBackWidget(titlew)
	layout.PushBackWidget(content)
	return layout
}

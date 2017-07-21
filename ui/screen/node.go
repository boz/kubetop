package screen

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const (
	nodeIndexPath = "/node"
	nodeShowPath  = "/node/show"
)

func RegisterNodeRoutes(router elements.Router) {
	router.Register(elements.NewRoute(nodeIndexPath), elements.NewHandler(nodeIndexHandler))
}

func NodeIndexRequest() elements.Request {
	return elements.NewRequest(nodeIndexPath)
}

type nodeIndex struct {
	content elements.Widget
	ctx     elements.Context
}

func nodeIndexHandler(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("node/index")

	db, err := ctx.Backend().Nodes()
	if err != nil {
		return nil, err
	}

	content := widget.NewNodeTable(ctx, db)
	index := &nodeIndex{content, ctx}
	content.Watch(index)

	return elements.NewScreen(ctx, req, "Nodes", index), nil
}

func (w *nodeIndex) Draw() {
	w.content.Draw()
}

func (w *nodeIndex) Resize() {
	w.content.Resize()
}

func (w *nodeIndex) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		w.Resize()
		return true
	}
	return w.content.HandleEvent(ev)
}

func (w *nodeIndex) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *nodeIndex) Size() (int, int) {
	return w.content.Size()
}

func (w *nodeIndex) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *nodeIndex) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}

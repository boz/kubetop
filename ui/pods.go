package ui

import (
	"fmt"

	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/gdamore/tcell/views"
)

type podIndexWidget struct {
	views.BoxLayout
	elements.Presentable
}

func newPodIndexWidget(p elements.Presenter) views.Widget {
	w := &podIndexWidget{
		BoxLayout: *views.NewBoxLayout(views.Vertical),
	}
	p.New("pods/index", w)
	go w.run()
	return w
}

func (w *podIndexWidget) run() {
	p := w.Presenter()

	ds, err := p.Backend().Pods()

	if err != nil {
		w.handleError(err, "datasource")
		return
	}

	pods, err := ds.List()
	if err != nil {
		w.handleError(err, "list")
		return
	}

	p.PostFunc(func() {
		w.initialize(pods)
	})

	sub := ds.Subscribe()
	defer sub.Close()

	for {
		select {
		case <-p.Closed():
			w.Env().Log().Debug("presenter closed")
			return
		case <-sub.Closed():
			w.Env().Log().Debug("sub closed")
			return
		case ev, ok := <-sub.Events():
			w.Env().Log().Debugf("got event %#v", ev)

			if !ok {
				w.Env().Log().Debug("events closed")
				// closed
				return
			}
			p.PostFunc(func() {
				w.handleDSEvent(ev)
			})
		}
	}
}

func (w *podIndexWidget) initialize(pods []pod.Pod) {
	for _, pod := range pods {
		txt := views.NewText()
		txt.SetText(pod.Name())
		w.InsertWidget(0, txt, 0.0)
		w.Resize()
	}
}

func (w *podIndexWidget) handleError(err error, msg string) {
	w.Env().LogErr(err, "pods")
	w.Presenter().PostFunc(func() {
		// show status error
		txt := views.NewText()
		txt.SetText(fmt.Sprint("ERROR: ", err))
		w.InsertWidget(0, txt, 0.0)
		w.Resize()
	})
}

func (w *podIndexWidget) handleDSEvent(ev pod.Event) {
	txt := views.NewText()
	txt.SetText(fmt.Sprintf("EVENT: %v %v", ev.Type(), ev.Resource()))
	w.InsertWidget(0, txt, 0.0)
	w.Resize()
}

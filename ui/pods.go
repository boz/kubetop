package ui

import (
	"github.com/boz/kcache"
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

/*

namespace name containers owner ctime phase condition message

- Labels
- Annotations
- CreationTimestamp
- OwnerReferences
- ResourceVersion
- ClusterName

- PodPhase
- Conditions
- Message
- Reason
- HostIP
- PodIP
- StartTime
- ContainerStatuses

- Volumes
- Containers
- RestartPolicy
- TerminationGracePeriodSeconds
- ActiveDeadlineSeconds
- DNSPolicy
- NodeSelector
*/

type podIndexWidget struct {
	content *elements.TableWidget
	model   elements.Table
	elements.Presentable
}

func newPodTable() elements.Table {
	header := elements.NewTableHeader([]elements.TableColumn{
		elements.NewTableTH("ns", "Namespace"),
		elements.NewTableTH("name", "Name"),
		elements.NewTableTH("version", "Version"),
		elements.NewTableTH("phase", "Phase"),
		elements.NewTableTH("message", "Message"),
	})
	rows := []elements.TableRow{}
	return elements.NewTable(header, rows)
}

func newPodIndexWidget(p elements.Presenter) views.Widget {

	model := newPodTable()

	w := &podIndexWidget{
		model:   model,
		content: elements.NewTableWidget(model),
	}

	p.New("pods/index", w)
	go w.run()
	return w
}

func (w *podIndexWidget) Draw() {
	w.content.Draw()
}

func (w *podIndexWidget) Resize() {
	w.content.Resize()
}

func (w *podIndexWidget) HandleEvent(ev tcell.Event) bool {
	return w.content.HandleEvent(ev)
}

func (w *podIndexWidget) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *podIndexWidget) Size() (int, int) {
	return w.content.Size()
}

func (w *podIndexWidget) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *podIndexWidget) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}

func (w *podIndexWidget) run() {
	p := w.Presenter()

	ds, err := p.Backend().Pods()
	if err != nil {
		w.handleError(err, "datasource")
		return
	}

	sub := ds.Subscribe()
	defer sub.Close()

	select {
	case <-p.Closed():
		w.Env().Log().Debug("presenter closed")
		return
	case <-sub.Closed():
		w.Env().Log().Debug("sub closed")
		return
	case <-sub.Ready():
		w.Env().Log().Debug("sub ready")
		pods, err := sub.List()
		if err != nil {
			w.handleError(err, "list")
			return
		}
		p.PostFunc(func() {
			w.initialize(pods)
		})
	}

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
		w.model.AddRow(w.rowForPod(pod))
	}
	w.Resize()
}

func (w *podIndexWidget) rowForPod(pod pod.Pod) elements.TableRow {

	stat := pod.Resource().Status

	phase := string(stat.Phase)
	message := stat.Message

	cols := []elements.TableColumn{
		elements.NewTableColumn("ns", pod.Resource().GetNamespace(), tcell.StyleDefault),
		elements.NewTableColumn("name", pod.Resource().GetName(), tcell.StyleDefault),
		elements.NewTableColumn("version", pod.Resource().GetResourceVersion(), tcell.StyleDefault),
		elements.NewTableColumn("phase", phase, tcell.StyleDefault),
		elements.NewTableColumn("message", message, tcell.StyleDefault),
	}
	return elements.NewTableRow(pod.ID(), cols)
}

func (w *podIndexWidget) handleError(err error, msg string) {
	w.Env().LogErr(err, "pods")
}

func (w *podIndexWidget) handleDSEvent(ev pod.Event) {
	switch ev.Type() {
	case kcache.EventTypeDelete:
		w.model.RemoveRow(ev.Resource().ID())
	case kcache.EventTypeCreate:
	case kcache.EventTypeUpdate:
		w.model.AddRow(w.rowForPod(ev.Resource()))
	}
}

package pod

import (
	"strings"

	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements/deflist"
	"github.com/boz/kubetop/ui/util"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"k8s.io/api/core/v1"
)

type Summary interface {
	pod.UnitaryHandler
	views.Widget
}

func NewSummary() Summary {
	return &summary{deflist.NewWidget(nil)}
}

type summary struct {
	leftdl deflist.Widget
}

func (w *summary) Draw() {
	w.leftdl.Draw()
}

func (w *summary) Resize() {
	w.leftdl.Resize()
}

func (w *summary) HandleEvent(ev tcell.Event) bool {
	return w.leftdl.HandleEvent(ev)
}

func (w *summary) SetView(view views.View) {
	w.leftdl.SetView(view)
}

func (w *summary) Size() (int, int) {
	return w.leftdl.Size()
}

func (w *summary) Watch(handler tcell.EventHandler) {
	w.leftdl.Watch(handler)
}

func (w *summary) Unwatch(handler tcell.EventHandler) {
	w.leftdl.Unwatch(handler)
}

func (w *summary) OnInitialize(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) OnCreate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) OnUpdate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) OnDelete(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) drawObject(obj *v1.Pod) {
	if obj == nil {
		return
	}

	var owners []string

	if len(obj.GetOwnerReferences()) == 0 {
		owners = append(owners, util.NASymbol)
	}

	for _, ref := range obj.GetOwnerReferences() {
		owners = append(owners, util.FormatOwnerReference(ref))
	}

	rows := []deflist.Row{
		deflist.NewSimpleRow("Name", obj.GetName()),
		deflist.NewSimpleRow("Namespace", obj.GetNamespace()),
		deflist.NewSimpleRow("Node", obj.Spec.NodeName),
		deflist.NewSimpleRow("Owners", strings.Join(owners, ",")),
	}

	w.leftdl.SetRows(rows)

	// obj.Status.StartTime

	/*
		text += "\n"

		nready := 0
		for _, cs := range obj.Status.ContainerStatuses {
			if cs.Ready {
				nready++
			}

		}

		text += fmt.Sprintf("containers ( %v/%v ready ):\n", nready, len(obj.Status.ContainerStatuses))

		text += "  name ready restarts state\n"

		for _, cs := range obj.Status.ContainerStatuses {
			text += fmt.Sprintf("  %v: %v %v ", cs.Name, cs.Ready, cs.RestartCount)

			switch {
			case cs.State.Waiting != nil:
				text += "W: " + cs.State.Waiting.Reason
			case cs.State.Running != nil:
				text += "R: " + cs.State.Running.StartedAt.String()
			case cs.State.Terminated != nil:
				text += "T: " + cs.State.Terminated.FinishedAt.String()
			default:
				text += "W: <unknown>"
			}
			text += "\n"
		}

		w.SetText(text)
	*/
}

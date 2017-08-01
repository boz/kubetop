package view

import (
	"fmt"
	"strings"

	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/elements/deflist"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/ui/util"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"k8s.io/api/core/v1"
)

func PodTableColumns() []table.TH {
	return []table.TH{
		table.NewTH("ns", "Namespace", true, 0),
		table.NewTH("name", "Name", true, 1),
		table.NewTH("status", "Status", true, -1),
		table.NewTH("containers", "Containers", true, -1),
		table.NewTH("message", "Message", true, -1),
	}
}

type podTable struct {
	content table.Display
}

func NewPodTableWriter(content table.Display) pod.Handler {
	return &podTable{content}
}

func (pt *podTable) OnInitialize(objs []*v1.Pod) {
	rows := make([]table.TR, 0, len(objs))
	for _, obj := range objs {
		rows = append(rows, pt.renderRow(obj))
	}
	pt.content.ResetRows(rows)
}

func (pt *podTable) OnCreate(obj *v1.Pod) {
	pt.content.InsertRow(pt.renderRow(obj))
}

func (pt *podTable) OnUpdate(obj *v1.Pod) {
	pt.content.UpdateRow(pt.renderRow(obj))
}

func (pt *podTable) OnDelete(obj *v1.Pod) {
	pt.content.RemoveRow(backend.ObjectID(obj))
}

func (pt *podTable) renderRow(obj *v1.Pod) table.TR {

	stat := obj.Status

	conditions := ""
	for _, c := range stat.Conditions {
		conditions += abbreviatePodConditionType(c.Type)
		switch c.Status {
		case v1.ConditionTrue:
			conditions += "+"
		case v1.ConditionFalse:
			conditions += "-"
		case v1.ConditionUnknown:
			conditions += "?"
		}
	}

	pstatus := conditions

	switch stat.Phase {
	case v1.PodPending:
		pstatus = "P " + conditions
	case v1.PodRunning:
		pstatus = "R " + conditions
	case v1.PodSucceeded:
		pstatus = "S " + conditions
	case v1.PodFailed:
		pstatus = "F " + conditions
	case v1.PodUnknown:
		pstatus = "U " + conditions
	}

	cready := 0
	cwaiting := 0
	crunning := 0
	cterminated := 0
	crestarts := int32(0)

	for _, cs := range obj.Status.ContainerStatuses {
		if cs.Ready {
			cready++
		}

		crestarts += cs.RestartCount

		switch {
		case cs.State.Waiting != nil && cs.State.Waiting.Reason != "":
			cwaiting++
		case cs.State.Running != nil:
			crunning++
		case cs.State.Terminated != nil:
			cterminated++
		default:
			cwaiting++
		}
	}

	cstatus := fmt.Sprintf("%v: %v/%v/%v (%v)",
		cready, cwaiting, crunning, cterminated, crestarts)

	cols := []table.TD{
		table.NewTD("ns", obj.GetNamespace(), theme.LabelNormal),
		table.NewTD("name", obj.GetName(), theme.LabelNormal),
		table.NewTD("status", pstatus, theme.LabelNormal),
		table.NewTD("containers", cstatus, theme.LabelNormal),
		table.NewTD("message", obj.Status.Message, theme.LabelNormal),
	}
	return table.NewTR(backend.ObjectID(obj), cols)
}

type PodSummary interface {
	pod.UnitaryHandler
	views.Widget
}

func NewPodSummary() PodSummary {
	return &podSummary{deflist.NewWidget(nil)}
}

type podSummary struct {
	leftdl deflist.Widget
}

func (w *podSummary) Draw() {
	w.leftdl.Draw()
}

func (w *podSummary) Resize() {
	w.leftdl.Resize()
}

func (w *podSummary) HandleEvent(ev tcell.Event) bool {
	return w.leftdl.HandleEvent(ev)
}

func (w *podSummary) SetView(view views.View) {
	w.leftdl.SetView(view)
}

func (w *podSummary) Size() (int, int) {
	return w.leftdl.Size()
}

func (w *podSummary) Watch(handler tcell.EventHandler) {
	w.leftdl.Watch(handler)
}

func (w *podSummary) Unwatch(handler tcell.EventHandler) {
	w.leftdl.Unwatch(handler)
}

func (w *podSummary) OnInitialize(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podSummary) OnCreate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podSummary) OnUpdate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podSummary) OnDelete(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podSummary) drawObject(obj *v1.Pod) {
	if obj == nil {
		return
	}

	var owners []string

	if len(obj.GetOwnerReferences()) == 0 {
		owners = append(owners, "N/A")
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
	/*

		text := "name: " + obj.GetName() + "\n"
		text += "namespace: " + obj.GetNamespace() + "\n"

		text += "owners: "

		for _, ref := range obj.GetOwnerReferences() {
			text += ref.Kind + "/" + ref.Name
		}

		text += "\n"
		text += "node: " + obj.Spec.NodeName + "\n"

		text += "start time: "
		if obj.Status.StartTime == nil {
			text += "N/A"
		} else {
			text += obj.Status.StartTime.String()
		}

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

func abbreviatePodConditionType(ct v1.PodConditionType) string {
	switch ct {
	case v1.PodScheduled:
		return "S"
	case v1.PodReady:
		return "R"
	case v1.PodInitialized:
		return "I"
	default:
		return "?"
	}
}

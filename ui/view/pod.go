package view

import (
	"fmt"

	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
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
		conditions += string(c.Type)[0:1]
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

type PodDetails interface {
	pod.UnitaryHandler
	views.Widget
}

func NewPodDetails() PodDetails {
	return &podDetails{*views.NewText()}
}

type podDetails struct {
	views.Text
}

func (w *podDetails) OnInitialize(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podDetails) OnCreate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podDetails) OnUpdate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podDetails) OnDelete(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *podDetails) drawObject(obj *v1.Pod) {
	if obj == nil {
		return
	}

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

	for _, cs := range obj.Status.ContainerStatuses {
		text += fmt.Sprintf("  %v: ready=%v restarts=%v state=", cs.Name, cs.Ready, cs.RestartCount)

		switch {
		case cs.State.Waiting != nil:
			text += "waiting: " + cs.State.Waiting.Reason
		case cs.State.Running != nil:
			text += "running: " + cs.State.Running.StartedAt.String()
		case cs.State.Terminated != nil:
			text += "terminated: " + cs.State.Terminated.FinishedAt.String()
		default:
			text += "waiting"
		}

	}

	text += "\n"

	w.SetText(text)
}

package view

import (
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/controller"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell/views"
	"k8s.io/api/core/v1"
)

func PodTableColumns() []table.TH {
	return []table.TH{
		table.NewTH("ns", "Namespace", true, 0),
		table.NewTH("name", "Name", true, 1),
		table.NewTH("version", "Version", true, -1),
		table.NewTH("phase", "Phase", true, -1),
		table.NewTH("conditions", "Conditions", true, -1),
		table.NewTH("message", "Message", true, -1),
	}
}

type podTable struct {
	content table.Display
}

func NewPodTableWriter(content table.Display) controller.PodsHandler {
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

	phase := string(stat.Phase)
	message := stat.Message

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

	cols := []table.TD{
		table.NewTD("ns", obj.GetNamespace(), theme.LabelNormal),
		table.NewTD("name", obj.GetName(), theme.LabelNormal),
		table.NewTD("version", obj.GetResourceVersion(), theme.LabelNormal),
		table.NewTD("phase", phase, theme.LabelNormal),
		table.NewTD("conditions", conditions, theme.LabelNormal),
		table.NewTD("message", message, theme.LabelNormal),
	}
	return table.NewTR(backend.ObjectID(obj), cols)
}

type PodDetails interface {
	controller.PodHandler
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

	w.SetText(text)
}

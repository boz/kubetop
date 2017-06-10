package view

import (
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/controller"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"k8s.io/client-go/pkg/api/v1"
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

func NewPodTableWriter(content table.Display) controller.PodHandler {
	return &podTable{content}
}

func (pt *podTable) OnInitialize(objs []pod.Pod) {
	for _, obj := range objs {
		pt.content.InsertRow(pt.renderRow(obj))
	}
}

func (pt *podTable) OnCreate(obj pod.Pod) {
	pt.content.InsertRow(pt.renderRow(obj))
}

func (pt *podTable) OnUpdate(obj pod.Pod) {
	pt.content.UpdateRow(pt.renderRow(obj))
}

func (pt *podTable) OnDelete(obj pod.Pod) {
	pt.content.RemoveRow(obj.ID())
}

func (pt *podTable) renderRow(obj pod.Pod) table.TR {

	stat := obj.Resource().Status

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
		table.NewTD("ns", obj.Resource().GetNamespace(), theme.LabelNormal),
		table.NewTD("name", obj.Resource().GetName(), theme.LabelNormal),
		table.NewTD("version", obj.Resource().GetResourceVersion(), theme.LabelNormal),
		table.NewTD("phase", phase, theme.LabelNormal),
		table.NewTD("conditions", conditions, theme.LabelNormal),
		table.NewTD("message", message, theme.LabelNormal),
	}
	return table.NewTR(obj.ID(), cols)
}

package pod

import (
	"fmt"

	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/ui/util"
	"k8s.io/api/core/v1"
)

func TableColumns() []table.TH {
	return []table.TH{
		table.NewTH("ns", "Namespace", true, 0),
		table.NewTH("name", "Name", true, 1),
		table.NewTH("status", "Status", true, -1),
		table.NewTH("containers", "Containers", true, -1),
		table.NewTH("age", "Age", true, -1),
		table.NewTH("message", "Message", true, -1),
	}
}

type podsTable struct {
	content table.Display
}

func NewTable(content table.Display) pod.Handler {
	return &podsTable{content}
}

func (pt *podsTable) OnInitialize(objs []*v1.Pod) {
	rows := make([]table.TR, 0, len(objs))
	for _, obj := range objs {
		rows = append(rows, pt.renderRow(obj))
	}
	pt.content.ResetRows(rows)
}

func (pt *podsTable) OnCreate(obj *v1.Pod) {
	pt.content.InsertRow(pt.renderRow(obj))
}

func (pt *podsTable) OnUpdate(obj *v1.Pod) {
	pt.content.UpdateRow(pt.renderRow(obj))
}

func (pt *podsTable) OnDelete(obj *v1.Pod) {
	pt.content.RemoveRow(backend.ObjectID(obj))
}

func (pt *podsTable) renderRow(obj *v1.Pod) table.TR {

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

	age := util.FormatAgeFromNow(obj.GetCreationTimestamp().Time)

	cols := []table.TD{
		table.NewTD("ns", obj.GetNamespace(), theme.LabelNormal),
		table.NewTD("name", obj.GetName(), theme.LabelNormal),
		table.NewTD("status", pstatus, theme.LabelNormal),
		table.NewTD("containers", cstatus, theme.LabelNormal),
		table.NewTD("age", age, theme.LabelNormal),
		table.NewTD("message", obj.Status.Message, theme.LabelNormal),
	}
	return table.NewTR(backend.ObjectID(obj), cols)
}

package pod

import (
	"fmt"

	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/ui/util"
	"k8s.io/api/core/v1"
)

func ContainersTableColumns() []table.TH {
	return []table.TH{
		table.NewTH("name", "Name", true, 0),
		table.NewTH("restarts", "Restarts", true, -1),
		table.NewTH("state", "State", true, -1),
	}
}

type containersTable struct {
	content table.Display
}

func NewContainersTable(content table.Display) pod.UnitaryHandler {
	return &containersTable{content}
}

func (t *containersTable) OnCreate(obj *v1.Pod) {
	t.resetRows(obj.Status.ContainerStatuses)
}

func (t *containersTable) OnUpdate(obj *v1.Pod) {
	t.resetRows(obj.Status.ContainerStatuses)
}

func (t *containersTable) OnDelete(obj *v1.Pod) {
	t.resetRows(nil)
}

func (t *containersTable) OnInitialize(obj *v1.Pod) {
	t.resetRows(obj.Status.ContainerStatuses)
}

func (t *containersTable) resetRows(objs []v1.ContainerStatus) {
	rows := make([]table.TR, 0, len(objs))
	for _, obj := range objs {
		rows = append(rows, t.renderRow(obj))
	}
	t.content.ResetRows(rows)
}

func (t *containersTable) renderRow(obj v1.ContainerStatus) table.TR {

	cols := []table.TD{
		table.NewTD("name", obj.Name, theme.LabelNormal),
		table.NewTD("restarts", fmt.Sprintf("%v", obj.RestartCount), theme.LabelNormal),
	}

	state := obj.State

	switch {
	case state.Waiting != nil:
		label := fmt.Sprintf("Waiting: %v",
			state.Waiting.Reason)
		cols = append(cols, table.NewTD("state", label, theme.LabelWarn))
	case state.Running != nil:
		label := fmt.Sprintf("Running: %v",
			util.FormatAgeFromNow(state.Running.StartedAt.Time))
		cols = append(cols, table.NewTD("state", label, theme.LabelNormal))
	case state.Terminated != nil:
		label := fmt.Sprintf("Terminated: %v",
			util.FormatAgeFromNow(state.Terminated.StartedAt.Time))
		cols = append(cols, table.NewTD("state", label, theme.LabelNormal))
	default:
		label := fmt.Sprintf("Waiting: unknown")
		cols = append(cols, table.NewTD("state", label, theme.LabelError))
	}
	return table.NewTR(obj.Name, cols)
}

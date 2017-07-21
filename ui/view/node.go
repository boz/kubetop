package view

import (
	"github.com/boz/kcache/types/node"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"k8s.io/api/core/v1"
)

func NodeTableColumns() []table.TH {
	return []table.TH{
		table.NewTH("name", "Name", true, 0),
		table.NewTH("phase", "Phase", true, -1),
	}
}

type nodeTable struct {
	content table.Display
}

func NewNodeTableWriter(content table.Display) node.Handler {
	return &nodeTable{content}
}

func (pt *nodeTable) OnInitialize(objs []*v1.Node) {
	rows := make([]table.TR, 0, len(objs))
	for _, obj := range objs {
		rows = append(rows, pt.renderRow(obj))
	}
	pt.content.ResetRows(rows)
}

func (pt *nodeTable) OnCreate(obj *v1.Node) {
	pt.content.InsertRow(pt.renderRow(obj))
}

func (pt *nodeTable) OnUpdate(obj *v1.Node) {
	pt.content.UpdateRow(pt.renderRow(obj))
}

func (pt *nodeTable) OnDelete(obj *v1.Node) {
	pt.content.RemoveRow(backend.ObjectID(obj))
}

func (pt *nodeTable) renderRow(obj *v1.Node) table.TR {
	cols := []table.TD{
		table.NewTD("name", obj.GetName(), theme.LabelNormal),
		table.NewTD("phase", string(obj.Status.Phase), theme.LabelNormal),
	}
	return table.NewTR(obj.GetName(), cols)
}

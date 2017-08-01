package service

import (
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"k8s.io/api/core/v1"
)

func TableColumns() []table.TH {
	return []table.TH{
		table.NewTH("ns", "Namespace", true, 0),
		table.NewTH("name", "Name", true, 1),
		table.NewTH("type", "Type", true, -1),
		table.NewTH("ip", "IP", true, -1),
	}
}

type servicesTable struct {
	content table.Display
}

func NewTable(content table.Display) service.Handler {
	return &servicesTable{content}
}

func (pt *servicesTable) OnInitialize(objs []*v1.Service) {
	rows := make([]table.TR, 0, len(objs))
	for _, obj := range objs {
		rows = append(rows, pt.renderRow(obj))
	}
	pt.content.ResetRows(rows)
}

func (pt *servicesTable) OnCreate(obj *v1.Service) {
	pt.content.InsertRow(pt.renderRow(obj))
}

func (pt *servicesTable) OnUpdate(obj *v1.Service) {
	pt.content.UpdateRow(pt.renderRow(obj))
}

func (pt *servicesTable) OnDelete(obj *v1.Service) {
	pt.content.RemoveRow(backend.ObjectID(obj))
}

func (pt *servicesTable) renderRow(obj *v1.Service) table.TR {
	cols := []table.TD{
		table.NewTD("ns", obj.GetNamespace(), theme.LabelNormal),
		table.NewTD("name", obj.GetName(), theme.LabelNormal),
		table.NewTD("type", string(obj.Spec.Type), theme.LabelNormal),
		table.NewTD("ip", obj.Spec.ClusterIP, theme.LabelNormal),
	}
	return table.NewTR(backend.ObjectID(obj), cols)
}

package view

import (
	"github.com/boz/kubetop/backend/service"
	"github.com/boz/kubetop/ui/controller"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
)

func ServiceTableColumns() []table.TH {
	return []table.TH{
		table.NewTH("ns", "Namespace", true, 0),
		table.NewTH("name", "Name", true, 1),
		table.NewTH("type", "Type", true, -1),
		table.NewTH("ip", "IP", true, -1),
	}
}

type serviceTable struct {
	content table.Display
}

func NewServiceTableWriter(content table.Display) controller.ServiceHandler {
	return &serviceTable{content}
}

func (pt *serviceTable) OnInitialize(objs []service.Service) {
	for _, obj := range objs {
		pt.content.InsertRow(pt.renderRow(obj))
	}
}

func (pt *serviceTable) OnCreate(obj service.Service) {
	pt.content.InsertRow(pt.renderRow(obj))
}

func (pt *serviceTable) OnUpdate(obj service.Service) {
	pt.content.UpdateRow(pt.renderRow(obj))
}

func (pt *serviceTable) OnDelete(obj service.Service) {
	pt.content.RemoveRow(obj.ID())
}

func (pt *serviceTable) renderRow(obj service.Service) table.TR {
	cols := []table.TD{
		table.NewTD("ns", obj.Resource().GetNamespace(), theme.LabelNormal),
		table.NewTD("name", obj.Resource().GetName(), theme.LabelNormal),
		table.NewTD("type", string(obj.Resource().Spec.Type), theme.LabelNormal),
		table.NewTD("ip", obj.Resource().Spec.ClusterIP, theme.LabelNormal),
	}
	return table.NewTR(obj.ID(), cols)
}

package table

import (
	"fmt"

	"github.com/boz/kcache/types/event"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"k8s.io/api/core/v1"
)

func TableColumns() []table.TH {
	return []table.TH{
		table.NewTH("timestamp", "Timestamp", true, 0),
		table.NewTH("type", "Type", true, -1),
		table.NewTH("reason", "Reason", true, -1),
		table.NewTH("object", "Object", true, -1),
		table.NewTH("message", "Message", true, -1),
	}
}

type eventsTable struct {
	content table.Display
}

func NewTable(content table.Display) event.Handler {
	return &eventsTable{content}
}

func (pt *eventsTable) OnInitialize(objs []*v1.Event) {
	rows := make([]table.TR, 0, len(objs))
	for _, obj := range objs {
		rows = append(rows, pt.renderRow(obj))
	}
	pt.content.ResetRows(rows)
}

func (pt *eventsTable) OnCreate(obj *v1.Event) {
	pt.content.InsertRow(pt.renderRow(obj))
}

func (pt *eventsTable) OnUpdate(obj *v1.Event) {
	pt.content.UpdateRow(pt.renderRow(obj))
}

func (pt *eventsTable) OnDelete(obj *v1.Event) {
	pt.content.RemoveRow(backend.ObjectID(obj))
}

func (pt *eventsTable) renderRow(obj *v1.Event) table.TR {

	object := eventFormatInvolvedObject(obj)

	cols := []table.TD{
		table.NewTD("timestamp", obj.GetCreationTimestamp().String(), theme.LabelNormal),
		table.NewTD("type", obj.Type, theme.LabelNormal),
		table.NewTD("reason", obj.Reason, theme.LabelNormal),
		table.NewTD("object", object, theme.LabelNormal),
		table.NewTD("message", obj.Message, theme.LabelNormal),
	}
	return table.NewTR(string(obj.UID), cols)
}

func eventFormatInvolvedObject(obj *v1.Event) string {

	io := obj.InvolvedObject

	switch io.Kind {
	case "Node":
		return fmt.Sprintf("Node{%s:%s}", io.Name, io.UID)
	default:
		return fmt.Sprintf("%s{%s/%s}", io.Kind, io.Namespace, io.Name)
	}

}

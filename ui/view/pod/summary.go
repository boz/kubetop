package pod

import (
	"strings"

	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/deflist"
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/ui/util"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"k8s.io/api/core/v1"
)

type Summary interface {
	pod.UnitaryHandler
	elements.Themeable
}

func NewSummary() Summary {
	return &summary{leftdl: deflist.NewWidget(nil)}
}

type summary struct {
	leftdl deflist.Widget
	theme  theme.Theme
}

func (w *summary) SetTheme(th theme.Theme) {
	w.theme = th
	w.leftdl.SetTheme(th)
}

func (w *summary) Theme() theme.Theme {
	return w.theme
}

func (w *summary) Draw() {
	w.leftdl.Draw()
}

func (w *summary) Resize() {
	w.leftdl.Resize()
}

func (w *summary) HandleEvent(ev tcell.Event) bool {
	return w.leftdl.HandleEvent(ev)
}

func (w *summary) SetView(view views.View) {
	w.leftdl.SetView(view)
}

func (w *summary) Size() (int, int) {
	return w.leftdl.Size()
}

func (w *summary) Watch(handler tcell.EventHandler) {
	w.leftdl.Watch(handler)
}

func (w *summary) Unwatch(handler tcell.EventHandler) {
	w.leftdl.Unwatch(handler)
}

func (w *summary) OnInitialize(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) OnCreate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) OnUpdate(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) OnDelete(obj *v1.Pod) {
	w.drawObject(obj)
}

func (w *summary) drawObject(obj *v1.Pod) {
	if obj == nil {
		return
	}

	var owners []string

	if len(obj.GetOwnerReferences()) == 0 {
		owners = append(owners, util.NASymbol)
	}

	for _, ref := range obj.GetOwnerReferences() {
		owners = append(owners, util.FormatOwnerReference(ref))
	}

	rows := []deflist.Row{
		deflist.NewSimpleRow("Name", obj.GetName(), theme.LabelNormal),
		deflist.NewSimpleRow("Namespace", obj.GetNamespace(), theme.LabelNormal),
		deflist.NewSimpleRow("Node", obj.Spec.NodeName, theme.LabelNormal),
		deflist.NewSimpleRow("Owners", strings.Join(owners, ","), theme.LabelNormal),
	}

	w.leftdl.SetRows(rows)
}

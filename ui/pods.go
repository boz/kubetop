package ui

import "github.com/gdamore/tcell/views"

type podIndexWidget struct {
	wbase
	views.BoxLayout
}

func newPodIndexWidget(base wbase) views.Widget {
	return &podIndexWidget{base, *views.NewBoxLayout(views.Vertical)}
}

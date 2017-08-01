package elements_test

import (
	"testing"

	"github.com/boz/kubetop/ui/elements"
	"github.com/gdamore/tcell/views"
	"github.com/stretchr/testify/assert"
)

func TestPanesPushBackWidget(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)

		p.PushBackWidget(w1)
		p.PushBackWidget(w2)
		p.PushBackWidget(w3)

		assert.Equal(t, []views.Widget{w1, w2, w3}, p.Widgets())
	}
}

func TestPanesPushFrontWidget(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)

		p.PushFrontWidget(w1)
		p.PushFrontWidget(w2)
		p.PushFrontWidget(w3)

		assert.Equal(t, []views.Widget{w3, w2, w1}, p.Widgets())
	}
}

func TestPanesRemoveWidget(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)
		p.PushBackWidget(w2)
		p.PushBackWidget(w3)

		p.RemoveWidget(w1)
		assert.Equal(t, []views.Widget{w2, w3}, p.Widgets())

		p.RemoveWidget(w1)
		assert.Equal(t, []views.Widget{w2, w3}, p.Widgets())

		for _, w := range p.Widgets() {
			p.RemoveWidget(w)
		}
		assert.Equal(t, []views.Widget{}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)
		p.PushBackWidget(w2)
		p.PushBackWidget(w3)

		p.RemoveWidget(w2)
		assert.Equal(t, []views.Widget{w1, w3}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)
		p.PushBackWidget(w2)
		p.PushBackWidget(w3)

		p.RemoveWidget(w3)
		assert.Equal(t, []views.Widget{w1, w2}, p.Widgets())
	}
}

func TestPanesInsertBeforeWidget(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)
		p.InsertBeforeWidget(w1, w1)

		assert.Empty(t, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)

		p.InsertBeforeWidget(w3, w2)
		assert.Equal(t, []views.Widget{w1}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)
		p.PushBackWidget(w3)

		p.InsertBeforeWidget(w1, w2)
		assert.Equal(t, []views.Widget{w2, w1, w3}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)
		p.PushBackWidget(w3)

		p.InsertBeforeWidget(w3, w2)
		assert.Equal(t, []views.Widget{w1, w2, w3}, p.Widgets())
	}
}

func TestPanesInsertAfterWidget(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)
		p.InsertAfterWidget(w1, w1)
		assert.Empty(t, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)

		p.InsertAfterWidget(w3, w2)
		assert.Equal(t, []views.Widget{w1}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)
		p.PushBackWidget(w3)

		p.InsertAfterWidget(w1, w2)
		assert.Equal(t, []views.Widget{w1, w2, w3}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.PushBackWidget(w1)
		p.PushBackWidget(w3)

		p.InsertAfterWidget(w3, w2)
		assert.Equal(t, []views.Widget{w1, w3, w2}, p.Widgets())
	}
}

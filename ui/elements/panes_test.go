package elements_test

import (
	"testing"

	"github.com/boz/kubetop/ui/elements"
	"github.com/gdamore/tcell/views"
	"github.com/stretchr/testify/assert"
)

func TestPanesAppend(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)

		p.Append(w1)
		p.Append(w2)
		p.Append(w3)

		assert.Equal(t, []views.Widget{w1, w2, w3}, p.Widgets())
	}
}

func TestPanesPrepend(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)

		p.Prepend(w1)
		p.Prepend(w2)
		p.Prepend(w3)

		assert.Equal(t, []views.Widget{w3, w2, w1}, p.Widgets())
	}
}

func TestPanesRemove(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)
		p.Append(w1)
		p.Append(w2)
		p.Append(w3)

		p.Remove(w1)
		assert.Equal(t, []views.Widget{w2, w3}, p.Widgets())

		p.Remove(w1)
		assert.Equal(t, []views.Widget{w2, w3}, p.Widgets())

		for _, w := range p.Widgets() {
			p.Remove(w)
		}
		assert.Equal(t, []views.Widget{}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)
		p.Append(w2)
		p.Append(w3)

		p.Remove(w2)
		assert.Equal(t, []views.Widget{w1, w3}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)
		p.Append(w2)
		p.Append(w3)

		p.Remove(w3)
		assert.Equal(t, []views.Widget{w1, w2}, p.Widgets())
	}
}

func TestPanesInsertBefore(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)
		p.InsertBefore(w1, w1)

		assert.Empty(t, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)

		p.InsertBefore(w3, w2)
		assert.Equal(t, []views.Widget{w1}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)
		p.Append(w3)

		p.InsertBefore(w1, w2)
		assert.Equal(t, []views.Widget{w2, w1, w3}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)
		p.Append(w3)

		p.InsertBefore(w3, w2)
		assert.Equal(t, []views.Widget{w1, w2, w3}, p.Widgets())
	}
}

func TestPanesInsertAfter(t *testing.T) {
	w1 := views.NewSpacer()
	w2 := views.NewSpacer()
	w3 := views.NewSpacer()

	{
		p := elements.NewVPanes(true)
		p.InsertAfter(w1, w1)
		assert.Empty(t, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)

		p.InsertAfter(w3, w2)
		assert.Equal(t, []views.Widget{w1}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)
		p.Append(w3)

		p.InsertAfter(w1, w2)
		assert.Equal(t, []views.Widget{w1, w2, w3}, p.Widgets())
	}

	{
		p := elements.NewVPanes(true)
		p.Append(w1)
		p.Append(w3)

		p.InsertAfter(w3, w2)
		assert.Equal(t, []views.Widget{w1, w3, w2}, p.Widgets())
	}
}

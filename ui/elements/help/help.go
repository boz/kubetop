package help

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/deflist"
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/util"
	"github.com/gdamore/tcell/views"
)

type Key interface {
	Label() string
	Description() string
}

func NewKey(label, description string) Key {
	return key{label, description}
}

type key struct {
	label       string
	description string
}

func (k key) Label() string {
	return k.label
}

func (k key) Description() string {
	return k.description
}

func NewSection(env util.Env, title string, keys []Key) views.Widget {

	outer := elements.NewVPanes(env, false)

	wtitle := views.NewText()
	wtitle.SetText(title)
	wtitle.SetAlignment(views.HAlignCenter)
	outer.Append(wtitle)

	if len(keys) == 0 {
		return outer
	}

	inner := elements.NewHPanes(env, true)

	var lkeys []deflist.Row
	var rkeys []deflist.Row

	for idx, key := range keys {
		row := deflist.NewSimpleRow(key.Label(), key.Description(), theme.LabelNormal)
		if idx%2 == 0 {
			lkeys = append(lkeys, row)
		} else {
			rkeys = append(rkeys, row)
		}
	}

	inner.Append(deflist.NewWidget(lkeys))
	if len(rkeys) > 0 {
		inner.Append(elements.AlignRight(deflist.NewWidget(rkeys)))
	}

	outer.Append(inner)

	return outer
}

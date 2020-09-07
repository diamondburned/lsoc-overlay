package reddot

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Dot struct {
	gtk.Revealer
	Dot *gtk.Label
}

func New(ms uint, markup string) *Dot {
	red, _ := gtk.LabelNew("")
	red.SetMarkup(markup)
	red.SetVAlign(gtk.ALIGN_START)
	red.Show()

	rctx, _ := red.GetStyleContext()
	rctx.AddClass("reddot")

	rev, _ := gtk.RevealerNew()
	rev.Add(red)
	rev.SetTransitionType(gtk.REVEALER_TRANSITION_TYPE_CROSSFADE)
	rev.SetTransitionDuration(ms)

	if ms > 0 {
		id, _ := glib.TimeoutAdd(ms, func() bool {
			rev.SetRevealChild(!rev.GetRevealChild())
			return true
		})

		// Stop the callback timer if the revealer is destroyed.
		rev.Connect("destroy", func() { glib.SourceRemove(id) })
	}

	return &Dot{
		Revealer: *rev,
		Dot:      red,
	}
}

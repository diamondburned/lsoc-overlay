package main

import (
	"flag"
	"image"
	"log"

	"github.com/diamondburned/lsoc-overlay/components/camerabox"
	"github.com/diamondburned/lsoc-overlay/components/reddot"
	"github.com/diamondburned/lsoc-overlay/gdkutil"
	"github.com/diamondburned/lsof"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"

	"net/http"
	_ "net/http/pprof"
)

var configPath = "./config.json"

func init() {
	flag.StringVar(&configPath, "c", configPath, "path to config.json")
	flag.Parse()

	gtk.Init(nil)

	// Load the application-specific CSS.
	LoadCSS(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION, `
		.background {
			background-color: transparent;
		}

		box.main {
			padding: 5px;
			border-radius: 5px;
			background-color: alpha(@theme_bg_color, 0.25);
		}

		@define-color recording #F04747;

		label.reddot {
			color: @recording;
			text-shadow: 0px 0px 2px alpha(@recording, 0.5);
		}
	`)

	go http.ListenAndServe(":20485", nil)
}

func main() {
	c, err := ReadConfig(configPath)
	if err != nil {
		log.Fatalln("Failed to read config:", err)
	}

	if err := c.LoadCSS(); err != nil {
		log.Fatalln("Failed to load CSS:", err)
	}

	if c.NumScanners > 0 {
		lsof.NumWorkers = c.NumScanners
	}

	red := reddot.New(c.RedBlinkMs, c.RedButton)
	red.Show()

	cam := camerabox.NewCameraBox()
	cam.SetHiddenProcs(c.HiddenProcs)
	cam.Show()

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	box.PackStart(red, false, false, 0)
	box.PackStart(cam, true, true, 5)
	box.Show()

	sctx, _ := box.GetStyleContext()
	sctx.AddClass("main")

	rev, _ := gtk.RevealerNew()
	rev.SetTransitionType(gtk.REVEALER_TRANSITION_TYPE_CROSSFADE)
	rev.SetTransitionDuration(100)
	rev.SetRevealChild(false)
	rev.Add(box)
	rev.Show()

	w, _ := gtk.WindowNew(gtk.WINDOW_POPUP)
	w.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
	w.SetSkipTaskbarHint(true)
	w.SetSkipPagerHint(true)
	w.SetEvents(0)
	w.Move(c.Window.X, c.Window.Y)
	setAlphaState(w)

	w.Add(rev)
	w.Show()

	if c.Window.Passthrough {
		setPassthrough(w)
	}

	glib.TimeoutAdd(c.PollingMs, func() bool {
		n, err := cam.Update()
		if err != nil {
			log.Println("Failed to update:", err)
		}

		// Reveal the overlay if there are cameras.
		if reveal := n > 0; reveal != rev.GetRevealChild() {
			rev.SetRevealChild(reveal)
		}

		return true
	})

	gtk.Main()
}

func getDefaultScreen() *gdk.Screen {
	d, _ := gdk.DisplayGetDefault()
	s, _ := d.GetDefaultScreen()
	return s
}

func LoadCSS(prio gtk.StyleProviderPriority, css string) error {
	prov, _ := gtk.CssProviderNew()

	if err := prov.LoadFromData(css); err != nil {
		return errors.Wrap(err, "Failed to parse CSS")
	}

	gtk.AddProviderForScreen(getDefaultScreen(), prov, uint(prio))
	return nil
}

type screener interface {
	GetScreen() *gdk.Screen
	SetVisual(v *gdk.Visual)
}

func setAlphaState(widget screener) {
	var screen = widget.GetScreen()

	var visual, _ = screen.GetRGBAVisual()
	// Fallback to the default system visual if there's no RGBA visual.
	if alpha := visual != nil; !alpha {
		visual, _ = screen.GetSystemVisual()
	}

	widget.SetVisual(visual)
}

type windower interface {
	GetWindow() (*gdk.Window, error)
}

func setPassthrough(widget windower) {
	w, _ := widget.GetWindow()
	gdkutil.WindowInputShapeCombineRegion(w, image.Rect(0, 0, 0, 0), 0, 0)
}

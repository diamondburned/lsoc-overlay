package camerabox

import (
	"fmt"
	"strings"

	"github.com/diamondburned/lsoc-overlay/camera"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

type CameraBox struct {
	gtk.Box

	hiddenProcs []string

	// states
	cameras map[string]*Camera
}

func NewCameraBox() *CameraBox {
	rev, _ := gtk.RevealerNew()
	rev.SetRevealChild(true)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	box.Show()

	return &CameraBox{
		Box:     *box,
		cameras: map[string]*Camera{},
	}
}

func (b *CameraBox) invalidate() {
	wipeChildren(b)
	b.cameras = map[string]*Camera{}
}

func (b *CameraBox) SetHiddenProcs(hiddenProcs []string) {
	b.hiddenProcs = hiddenProcs
}

// Update returns the current number of visible cameras after the update as well
// as an error if any.
func (b *CameraBox) Update() (n int, err error) {
	c, err := camera.Cameras()
	if err != nil {
		b.invalidate()
		return 0, errors.Wrap(err, "Failed to get cameras")
	}

	c = camera.FilterCameras(c, func(c *camera.Camera) bool {
		// Filter out inactive cameras.
		if !c.IsActive() {
			return false
		}

		// Filter out cameras with ignored process names.
		for _, procName := range b.hiddenProcs {
			for _, proc := range c.PIDs {
				var filtered = c.PIDs[:0]

				if proc.Executable() == procName {
					// Ignore the camera entirely if this is the only process.
					if len(c.PIDs) == 1 {
						return false
					}
				} else {
					// Else, filter it out.
					filtered = append(filtered, proc)
				}

				// Set the filtered slice into the PIDs field of the camera
				// pointer.
				c.PIDs = filtered
			}
		}

		return true
	})

	// Drop whatever isn't in the camera list anymore.
Loop:
	for path, oldc := range b.cameras {
		for _, cam := range c {
			if cam.Path == path {
				continue Loop
			}
		}

		b.Box.Remove(oldc)
		delete(b.cameras, path)
	}

	// Guarantee that labels are added.
	for _, cam := range c {
		cm, ok := b.cameras[cam.Path]
		if !ok {
			cm = NewCamera()
			cm.Show()

			b.Box.Add(cm)
			b.cameras[cam.Path] = cm
		}

		// Update the label.
		cm.Update(cam)
	}

	return len(c), nil
}

type Camera struct {
	gtk.Box
	Name  *gtk.Label
	Procs *gtk.Label

	// states
	name  string
	execs []string
}

func NewCamera() *Camera {
	name, _ := gtk.LabelNew("")
	name.SetXAlign(0)
	name.SetYAlign(1)
	name.Show()

	procs, _ := gtk.LabelNew("")
	procs.SetXAlign(0)
	procs.SetYAlign(1)
	procs.Show()

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	box.Add(name)
	box.Add(procs)

	return &Camera{
		Box:   *box,
		Name:  name,
		Procs: procs,
	}
}

func (c *Camera) Update(cam camera.Camera) {
	if name := cam.Name(); c.name != name {
		c.name = name
		c.Name.SetLabel(cam.Name() + ": ")
	}

	if xs := camera.ExecutableNames(cam.PIDs); !stringsEq(c.execs, xs) {
		c.execs = xs
		c.Procs.SetMarkup(fmt.Sprintf(`<span size="small">%s</span>`, strings.Join(xs, ", ")))
	}
}

func stringsEq(i1, i2 []string) bool {
	if len(i1) != len(i2) {
		return false
	}

	for i := range i1 {
		if i1[i] != i2[i] {
			return false
		}
	}

	return true
}

type container interface {
	GetChildren() *glib.List
	Remove(gtk.IWidget)
}

func wipeChildren(w container) {
	w.GetChildren().Foreach(func(v interface{}) {
		w.Remove(v.(gtk.IWidget))
	})
}

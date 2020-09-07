package gdkutil

// #cgo pkg-config: cairo gdk-3.0
// #include "gdk.h"
import "C"

import (
	"image"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
)

func WindowInputShapeCombineRegion(window *gdk.Window, rect image.Rectangle, offsetX, offsetY int) {
	C.window_input_shape_combine_region(
		(*C.GdkWindow)(unsafe.Pointer(window.Native())),
		C.int(rect.Min.X),
		C.int(rect.Min.Y),
		C.int(rect.Dx()),
		C.int(rect.Dy()),
		C.int(offsetX),
		C.int(offsetY),
	)
}

#include <cairo/cairo.h>
#include <gdk/gdk.h>

void window_input_shape_combine_region(
	GdkWindow *window,
	int x, int y, int w, int h,
	int offset_x, int offset_y
) {

	cairo_rectangle_int_t rect = { x, y, w, h };
	cairo_region_t* region = cairo_region_create_rectangle(&rect);

	gdk_window_input_shape_combine_region(window, region, offset_x, offset_y);

	cairo_region_destroy(region);
};

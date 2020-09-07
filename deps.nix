pkgs: {
	buildInputs = with pkgs.gnome3; [
		glib gtk
	];
	nativeBuildInputs = with pkgs; [
		pkgconfig
	];
}

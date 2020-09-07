# lsoc-overlay

Overlay for the List of Open Cameras.

![screenshot](scrot.png)

## Install

Dependencies: `go`, `gtk3`, `v4l2`, Linux.

```sh
# Optionally add this to shellrc.
export PATH="$PATH:${GOBIN:-$GOPATH/bin}"

go install github.com/diamondburned/lsoc-overlay
```

## Usage

```sh
# Optionally install the default config.json.
cp config.json ~/.config/lsoc.config.json

lsoc-overlay [-c ~/.config/lsoc.config.json]
```

By default, `lsoc-overlay` will scan with a single thread every 1.2s to save
CPU time and battery. If one wishes to prioritize responsiveness, change
`num_scanners` to `-1` and `polling_ms` to `250` (or any low number). This will
ensure that `lsoc-overlay` scans with multiple threads in a short amount of
time.

### Nix

```nix
{ config, pkgs, lib, ... }:

let lsoc-overlay = builtins.fetchGit {
	url = "https://github.com/diamondburned/lsoc-overlay.git";
	rev = "1e56cc825ce885a64a7e635f54ff704958e7d029"; # update this
};

in {
	# <username> should be changed to the appropriate username.
	home-manager.users.<username> = {
		imports = [ "${lsoc-overlay}" ];
	
		services.lsoc-overlay = {
			enable = true;
			config = {
				polling_ms = 1000;
				num_scanners = -1;
			};
		};
	};
}
```

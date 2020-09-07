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

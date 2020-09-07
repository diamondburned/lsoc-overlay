package camera

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/diamondburned/lsof"
	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"github.com/reiver/go-v4l2"
)

func ExecutableNames(processes []ps.Process) []string {
	var names = make([]string, len(processes))
	for i, proc := range processes {
		names[i] = proc.Executable()
	}
	return names
}

func FilterCameras(cameras []Camera, filterFn func(c *Camera) bool) []Camera {
	var filtered = cameras[:0]

	for _, cam := range cameras {
		if filterFn(&cam) {
			filtered = append(filtered, cam)
		}
	}

	return filtered
}

type Camera struct {
	Path string // /dev/videoX
	PIDs []ps.Process

	name string
}

var ErrCameraNotFound = errors.New("camera not found")

func OpenCamera(path string) (*Camera, error) {
	c, err := cameras(path)
	if err != nil {
		return nil, err
	}

	if len(c) == 0 {
		return nil, ErrCameraNotFound
	}

	return &c[0], nil
}

func Cameras() ([]Camera, error) {
	return cameras("/dev/video")
}

func cameras(prefix string) ([]Camera, error) {
	l, err := lsof.Scan(func(str string) bool { return strings.HasPrefix(str, prefix) })
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get cameras")
	}

	var cams = make([]Camera, 0, len(l))

	for file, pids := range l {
		// Sort the PIDs.
		sort.Ints(pids)

		var cam = Camera{
			Path: file,
			PIDs: make([]ps.Process, 0, len(pids)),
		}

		for _, pid := range pids {
			p, err := ps.FindProcess(pid)
			if err != nil {
				log.Println("Failed to find process with PID", pid)
				continue
			}

			cam.PIDs = append(cam.PIDs, p)
		}

		cams = append(cams, cam)
	}

	sort.Slice(cams, func(i, j int) bool {
		return cams[i].Path < cams[j].Path
	})

	return cams, nil
}

func (c Camera) IsActive() bool {
	return len(c.PIDs) > 0
}

func (c Camera) Name() string {
	f, err := v4l2.Open(c.Path)
	if err != nil {
		log.Println("Failed to open webcam:", err)
		return fmt.Sprintf("Unknown webcam (%s)", c.Path)
	}

	defer f.Close()

	s, err := f.Card()
	if err != nil {
		log.Println("Failed to read webcam name:", err)
		return fmt.Sprintf("Unknown webcam (%s)", c.Path)
	}

	return s
}

func (c Camera) Configure(fn func(camera *v4l2.Device) error) error {
	f, err := v4l2.Open(c.Path)
	if err != nil {
		return errors.Wrap(err, "Failed to open webcam")
	}

	defer f.Close()

	fn(f)

	return nil
}

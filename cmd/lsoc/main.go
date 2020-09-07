package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/diamondburned/lsoc-overlay/camera"
)

var JSON bool

func init() {
	flag.BoolVar(&JSON, "j", JSON, "JSON output")
	flag.Parse()

	log.SetFlags(0)
}

func main() {
	c, err := camera.Cameras()
	if err != nil {
		log.Fatalln("Failed to get cameras:", err)
	}

	if !JSON {
		for _, cam := range c {
			fmt.Printf("%s (%s)\n", cam.Name(), cam.Path)

			for _, proc := range cam.Procs {
				fmt.Printf("\t- %s (%d)\n", proc.Executable(), proc.Pid())
			}
		}
	} else {
		if err := json.NewEncoder(os.Stdout).Encode(NewJSONCameras(c)); err != nil {
			log.Fatalln("Failed to encode cameras into JSON:", err)
		}
	}
}

type JSONCamera struct {
	Name string `json:"name"`
	Path string `json:"path"`
	PIDs []int  `json:"pids"`
}

func NewJSONCameras(cameras []camera.Camera) []JSONCamera {
	var j = make([]JSONCamera, len(cameras))
	for i, camera := range cameras {
		j[i] = NewJSONCamera(camera)
	}
	return j
}

func NewJSONCamera(c camera.Camera) JSONCamera {
	var pids = make([]int, len(c.Procs))
	for i, proc := range c.Procs {
		pids[i] = proc.Pid()
	}

	return JSONCamera{
		Name: c.Name(),
		Path: c.Path,
		PIDs: pids,
	}
}

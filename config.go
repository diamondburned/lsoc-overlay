package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

type Config struct {
	RedButton  string `json:"red_button"`
	RedBlinkMs uint   `json:"red_blink_ms"`

	CustomCSS string `json:"custom_css"`
	PollingMs uint   `json:"polling_ms"`

	HiddenProcs []string `json:"hidden_procs"`
	NumScanners int      `json:"num_scanners"`

	Window struct {
		X           int  `json:"x"`
		Y           int  `json:"y"`
		Passthrough bool `json:"passthrough"`
	} `json:"window"`
}

func ReadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open config")
	}

	defer f.Close()

	var c = Config{
		RedButton: "⬤",
		PollingMs: 500,
	}

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, errors.Wrap(err, "Failed to decode JSON")
	}

	return &c, nil
}

func (c *Config) LoadCSS() error {
	if c.CustomCSS == "" {
		return nil
	}

	b, err := ioutil.ReadFile(c.CustomCSS)
	if err != nil {
		return err
	}

	return LoadCSS(gtk.STYLE_PROVIDER_PRIORITY_USER, string(b))
}

package main

import (
	"encoding/json"
	"os/exec"
)

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type OutputMode struct {
	Width   int `json:"width"`
	Height  int `json:"height"`
	Refresh int `json:"refresh"`
}

type Output struct {
	Id                 int              `json:"id"`
	Name               string           `json:"name"`
	Rect               Rect         `json:"rect"`
	Focus              []int            `json:"focus"`
	Border             string           `json:"border"`
	CurrentBorderWidth int              `json:"current_border_width"`
	Layout             string           `json:"layout"`
	Orientation        string           `json:"orientation"`
	Percent            float64          `json:"percent"`
	WindowRect         Rect         `json:"window_rect"`
	DecoRect           Rect         `json:"deco_rect"`
	Geometry           Rect         `json:"geometry"`
	Window             interface{}      `json:"window"`
	Urgent             bool             `json:"urgent"`
	FloatingNodes      []interface{}    `json:"floating_nodes"`
	Sticky             bool             `json:"sticky"`
	Type               string           `json:"type"`
	Active             bool             `json:"active"`
	Primary            bool             `json:"primary"`
	Make               string           `json:"make"`
	Model              string           `json:"model"`
	Serial             string           `json:"serial"`
	Scale              float64          `json:"scale"`
	Transform          string           `json:"transform"`
	CurrentWorkspace   string           `json:"current_workspace"`
	Modes              []OutputMode `json:"modes"`
	CurrentMode        OutputMode   `json:"current_mode"`
	Focused            bool             `json:"focused"`
	SubpixelHinting    string           `json:"subpixel_hinting"`
}

type Input struct {
	Identifier string            `json:"identifier"`
	Name       string            `json:"name"`
	Vendor     int               `json:"vendor"`
	Product    int               `json:"product"`
	Type       string            `json:"type"`
	Libinput   map[string]string `json:"libinput"`
}

func SwayMsg(outjson interface{}, args ...string) (err error) {
	var stdout []byte
	stdout, err = exec.Command("swaymsg", args...).Output()

	if err == nil && stdout != nil && outjson != nil {
		return json.Unmarshal(stdout, outjson)
	}

	return
}

func GetOutputs() (outputs []Output, err error) {
	err = SwayMsg(&outputs, "-t", "get_outputs")
	return
}

func GetOutputNames() (names []string, err error) {
	var outputs []Output
	outputs, err = GetOutputs()

	if err != nil {
		return
	}

	names = make([]string, len(outputs))

	for i, output := range outputs {
		names[i] = output.Name
	}

	return
}

func GetInputs() (outputs []Input, err error) {
	err = SwayMsg(&outputs, "-t", "get_inputs")
	return
}

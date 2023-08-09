package controller

import (
	"github.com.bisoncorp.autostrade/gameapi"
	"image/color"
	"sync"
)

type ColorChangedCallback func(color.Color)

type ColorableController struct {
	ColorableObject gameapi.Colorable
	callbacks       []ColorChangedCallback
	callbacksMu     sync.Mutex
}

func NewColorableController(colorableObject gameapi.Colorable) *ColorableController {
	return &ColorableController{ColorableObject: colorableObject, callbacks: make([]ColorChangedCallback, 0)}
}

func (c *ColorableController) AddCallback(fn ColorChangedCallback) {
	c.callbacksMu.Lock()
	defer c.callbacksMu.Unlock()
	c.callbacks = append(c.callbacks, fn)
}

func (c *ColorableController) Color() color.Color {
	c.callbacksMu.Lock()
	defer c.callbacksMu.Unlock()
	return c.ColorableObject.Color()
}

func (c *ColorableController) SetColor(col color.Color) {
	c.callbacksMu.Lock()
	defer c.callbacksMu.Unlock()
	c.ColorableObject.SetColor(col)
	c.callAll(col)
}

func (c *ColorableController) callAll(col color.Color) {
	for _, fn := range c.callbacks {
		if fn != nil {
			fn(col)
		}
	}
}

type ColorableBuffer struct {
	col color.Color
}

func NewColorableBuffer(col color.Color) *ColorableBuffer {
	return &ColorableBuffer{col: col}
}

func (c *ColorableBuffer) Color() color.Color {
	return c.col
}

func (c *ColorableBuffer) SetColor(c2 color.Color) {
	c.col = c2
}

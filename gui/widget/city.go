package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	api "github.com/bisoncorp/autostrade/gameapi"
	"image/color"
)

type City struct {
	widget.BaseWidget
	onTapped func(data api.CityData)
	data     api.CityData
	hover    bool
}

func (c *City) SetData(data api.CityData) {
	c.data = data
	c.Refresh()
}

func (c *City) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}
func (c *City) MouseIn(_ *desktop.MouseEvent) {
	c.hover = true
	c.Refresh()
}
func (c *City) MouseMoved(_ *desktop.MouseEvent) {}
func (c *City) MouseOut() {
	c.hover = false
	c.Refresh()
}
func (c *City) Tapped(e *fyne.PointEvent) {
	c.onTapped(c.data)
}
func (c *City) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(cityDimension, cityDimension))

	hoverRect := canvas.NewRectangle(color.Transparent)
	hoverRect.SetMinSize(fyne.NewSize(cityDimension+theme.Padding(), cityDimension+theme.Padding()))

	objects := []fyne.CanvasObject{hoverRect, rect}
	return &cityRenderer{
		wid:       c,
		objects:   objects,
		rect:      rect,
		hoverRect: hoverRect,
	}
}

func NewCity(data api.CityData, onTapped func(api.CityData)) *City {
	c := &City{data: data, onTapped: onTapped}
	c.ExtendBaseWidget(c)
	c.Refresh()
	return c
}

type cityRenderer struct {
	wid       *City
	objects   []fyne.CanvasObject
	rect      *canvas.Rectangle
	hoverRect *canvas.Rectangle
}

func (c *cityRenderer) Destroy() {}

func (c *cityRenderer) Layout(size fyne.Size) {
	c.hoverRect.Resize(size)
	c.rect.Resize(size.SubtractWidthHeight(theme.Padding(), theme.Padding()))
	c.rect.Move(fyne.NewPos(theme.Padding()/2, theme.Padding()/2))
}

func (c *cityRenderer) MinSize() fyne.Size {
	return c.hoverRect.MinSize()
}

func (c *cityRenderer) Objects() []fyne.CanvasObject {
	return c.objects
}

func (c *cityRenderer) Refresh() {
	if c.wid.hover {
		c.hoverRect.Show()
	} else {
		c.hoverRect.Hide()
	}
	data := c.wid.data
	c.rect.FillColor = data.Color
	c.rect.Refresh()
	c.hoverRect.FillColor = hoverColor(data.Color)
	c.hoverRect.Refresh()
}

func hoverColor(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 128,
	}
}

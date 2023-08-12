package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
	"image/color"
)

func vehicleMinSize() fyne.Size {
	return fyne.NewSize(vehicleDimension+theme.Padding(), vehicleDimension+theme.Padding())
}

type Vehicle struct {
	widget.BaseWidget

	hook     api.Vehicle
	onTapped func(api.Vehicle)
	hover    bool
}

func (v *Vehicle) center() fyne.Position {
	size := v.Size()
	return fyne.NewPos(size.Width/2, size.Height/2)
}

func (v *Vehicle) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}
func (v *Vehicle) MouseIn(_ *desktop.MouseEvent) {
	v.hover = true
	v.Refresh()
}
func (v *Vehicle) MouseMoved(_ *desktop.MouseEvent) {}
func (v *Vehicle) MouseOut() {
	v.hover = false
	v.Refresh()
}
func (v *Vehicle) Tapped(_ *fyne.PointEvent) {
	v.onTapped(v.hook)
}

func (v *Vehicle) CreateRenderer() fyne.WidgetRenderer {
	circle := canvas.NewCircle(color.Transparent)
	hoverCircle := canvas.NewCircle(color.Transparent)
	return &vehicleRenderer{
		wid:         v,
		objects:     []fyne.CanvasObject{hoverCircle, circle},
		circle:      circle,
		hoverCircle: hoverCircle,
	}
}

func NewVehicle(hook api.Vehicle, onTapped func(api.Vehicle)) *Vehicle {
	v := &Vehicle{hook: hook, onTapped: onTapped}
	v.ExtendBaseWidget(v)
	v.Refresh()
	return v
}

type vehicleRenderer struct {
	wid                 *Vehicle
	objects             []fyne.CanvasObject
	circle, hoverCircle *canvas.Circle
}

func (v *vehicleRenderer) Destroy() {}
func (v *vehicleRenderer) Layout(size fyne.Size) {
	v.hoverCircle.Resize(size)
	v.circle.Resize(size.SubtractWidthHeight(theme.Padding(), theme.Padding()))
	v.circle.Move(fyne.NewPos(theme.Padding()/2, theme.Padding()/2))
}
func (v *vehicleRenderer) MinSize() fyne.Size {
	return vehicleMinSize()
}

func (v *vehicleRenderer) Objects() []fyne.CanvasObject {
	return v.objects
}
func (v *vehicleRenderer) Refresh() {
	if v.wid.hover {
		v.hoverCircle.Show()
	} else {
		v.hoverCircle.Hide()
	}
	col := v.wid.hook.Color()
	v.circle.FillColor = col
	v.circle.Refresh()
	v.hoverCircle.FillColor = hoverColor(col)
	v.hoverCircle.Refresh()
}

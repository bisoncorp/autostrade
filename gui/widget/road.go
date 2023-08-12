package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
	"math"
)

type Road struct {
	widget.BaseWidget

	hook      api.Road
	highlight bool
}

func (r *Road) CreateRenderer() fyne.WidgetRenderer {
	line := canvas.NewLine(theme.ForegroundColor())
	line.StrokeWidth = roadDimension
	return &roadRenderer{
		wid:     r,
		objects: []fyne.CanvasObject{line},
		line:    line,
	}
}

func NewRoad(hook api.Road) *Road {
	r := &Road{hook: hook}
	r.ExtendBaseWidget(r)
	r.Refresh()
	return r
}

type roadRenderer struct {
	wid     *Road
	objects []fyne.CanvasObject
	line    *canvas.Line
}

func (r *roadRenderer) Destroy() {}
func (r *roadRenderer) Layout(size fyne.Size) {
	r.line.Resize(size)
}
func (r *roadRenderer) MinSize() fyne.Size {
	srcPos, dstPos := r.wid.hook.Src().Position(), r.wid.hook.Dst().Position()
	w := math.Abs(srcPos.X - dstPos.X)
	h := math.Abs(srcPos.Y - dstPos.Y)
	return fyne.NewSize(float32(w), float32(h))
}
func (r *roadRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
func (r *roadRenderer) Refresh() {
	if r.wid.highlight {
		r.line.StrokeColor = theme.PrimaryColor()
	} else {
		r.line.StrokeColor = theme.ForegroundColor()
	}
}

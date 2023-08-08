package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
)

type Road struct {
	widget.BaseWidget
	data     api.RoadData
	src, dst api.CityData
}

func (r *Road) SetData(data api.RoadData, src, dst api.CityData) {
	r.data, r.src, r.dst = data, src, dst
	r.Refresh()
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

func NewRoad() *Road {
	r := &Road{}
	r.ExtendBaseWidget(r)
	return r
}

type roadRenderer struct {
	wid     *Road
	objects []fyne.CanvasObject
	line    *canvas.Line
}

func (r *roadRenderer) Destroy() {}
func (r *roadRenderer) Layout(_ fyne.Size) {
	r.line.Position1 = scale(r.wid.src.Pos.ToPos32(), scaleFactor)
	r.line.Position2 = scale(r.wid.dst.Pos.ToPos32(), scaleFactor)
	r.line.Refresh()
}
func (r *roadRenderer) MinSize() fyne.Size {
	return r.line.MinSize()
}
func (r *roadRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
func (r *roadRenderer) Refresh() {
	r.Layout(fyne.Size{})
}

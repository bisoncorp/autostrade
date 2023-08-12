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

	line *canvas.Line
	hook api.Road
}

func (r *Road) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.line)
}

func NewRoad(hook api.Road) *Road {
	line := canvas.NewLine(theme.ForegroundColor())
	line.StrokeWidth = roadDimension
	r := &Road{hook: hook, line: line}
	r.ExtendBaseWidget(r)
	r.Refresh()
	return r
}

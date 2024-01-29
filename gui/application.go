package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	api "github.com/bisoncorp/autostrade/gameapi"
)

type Application struct {
	appl fyne.App
}

func NewApplication() *Application {
	a := &Application{}
	a.appl = app.NewWithID("github.com.bisoncorp.autostrade")
	return a
}

func (a *Application) Run() {
	a.appl.Run()
}

func (a *Application) NewWindow(sim api.Simulation) {
	wind := a.appl.NewWindow("Simulation")
	ui, menu := buildSimulationUi(sim, wind, a)
	wind.SetContent(ui)
	wind.SetMainMenu(menu)
	wind.Show()
}

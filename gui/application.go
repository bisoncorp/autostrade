package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com.bisoncorp.autostrade/game"
	api "github.com.bisoncorp.autostrade/gameapi"
)

type Application struct {
	appl       fyne.App
	wind       fyne.Window
	simulation api.Simulation
}

func NewApplication() *Application {
	a := &Application{}
	a.appl = app.NewWithID("github.com.bisoncorp.autostrade")
	a.wind = a.appl.NewWindow("Autostrade")
	a.simulation = game.NewFromData(api.ReadSimulationData("schema.json"))

	ui, menu := buildSimulationUi(a.simulation, a.wind)
	a.wind.SetContent(ui)
	a.wind.SetMainMenu(menu)

	return a
}

func (a *Application) ShowAndRun() {
	a.wind.ShowAndRun()
}

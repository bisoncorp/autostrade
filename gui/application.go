package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com.bisoncorp.autostrade/game"
	api "github.com.bisoncorp.autostrade/gameapi"
	gamewid "github.com.bisoncorp.autostrade/gui/widget"
	"image/color"
	"log"
	"time"
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

	m := gamewid.NewMap(2)
	m.OnCityTapped = func(data api.CityData) {
		log.Println(data)
	}
	m.OnVehicleTapped = func(data api.VehicleData) {
		log.Println(data)
	}
	a.wind.SetContent(m)
	a.simulation.Start()

	go func() {
		t := time.Tick(time.Second / 60)
		for {
			<-t
			m.SetData(a.simulation.PackData())
		}
	}()
	return a
}

func (a *Application) ShowAndRun() {
	a.wind.ShowAndRun()
}

func colorToRgba(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

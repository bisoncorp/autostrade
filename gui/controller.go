package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
	gamewid "github.com.bisoncorp.autostrade/gui/widget"
)

type speedString float64

func (s speedString) String() string {
	return fmt.Sprint("Simulation speed: ", int(s))
}

type simulationPlayController struct {
	mainWindow                      fyne.Window
	simulation                      api.Simulation
	mapWidget                       *gamewid.Map
	startButton, stopButton         *widget.Button
	startMenuButton, stopMenuButton *fyne.MenuItem
	speedLabel                      *widget.Label
	speedSlider                     *widget.Slider
}

func (s *simulationPlayController) startSimulation() {
	s.mapWidget.Start()
	s.startButton.Disable()
	s.startMenuButton.Disabled = true
	s.simulation.Start()
	s.stopButton.Enable()
	s.stopMenuButton.Disabled = false
}
func (s *simulationPlayController) stopSimulation() {
	s.stopMenuButton.Disabled = true
	s.stopButton.Disable()
	s.simulation.Stop()
	s.startMenuButton.Disabled = false
	s.startButton.Enable()
	s.mapWidget.Stop()
}
func (s *simulationPlayController) changeSpeedSimulation(value float64) {
	s.speedLabel.SetText(speedString(value).String())
	s.speedSlider.Value = value
	s.speedSlider.Refresh()
	s.simulation.SetSpeed(value)
}

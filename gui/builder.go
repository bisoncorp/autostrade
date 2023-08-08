package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
	gamewid "github.com.bisoncorp.autostrade/gui/widget"
	"image/color"
	"time"
)

func buildSimulationUi(sim api.Simulation, window fyne.Window) (fyne.CanvasObject, *fyne.MainMenu) {
	leftCnt := container.NewVBox()
	rightCnt := container.NewVBox()

	mapWidget := gamewid.NewMap()
	mapWidget.OnCityTapped = func(data api.CityData) {
		c := sim.City(data.Name)
		leftCnt.Add(gamewid.NewMiniWindow(fmt.Sprintf("City Property [%s]", data.Name), buildCityForm(c, window), func(w *gamewid.MiniWindow) {
			leftCnt.Remove(w)
		}))
	}
	mapWidget.OnVehicleTapped = func(data api.VehicleData) {
		v := sim.Vehicle(data.Plate)
		content, closeFn := buildVehicleForm(v, window)
		rightCnt.Add(gamewid.NewMiniWindow(fmt.Sprintf("Vehicle Property [%s]", data.Plate), content, func(w *gamewid.MiniWindow) {
			rightCnt.Remove(w)
			closeFn()
		}))
	}
	mapWidget.OnDataRequired = func() api.SimulationData {
		return sim.PackData()
	}

	playCtrl := &simulationPlayController{
		mainWindow: window,
		simulation: sim,
		mapWidget:  mapWidget,
	}

	scrollMap := container.NewScroll(mapWidget)
	scrollMap.SetMinSize(fyne.NewSize(300, 300))
	return container.NewBorder(nil, buildSimulationPlayControlBar(sim, playCtrl), leftCnt, rightCnt, scrollMap), buildMenu(playCtrl)
}

func buildSimulationPlayControlBar(sim api.Simulation, controller *simulationPlayController) fyne.CanvasObject {
	startButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), controller.startSimulation)
	startButton.Importance = widget.LowImportance

	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), controller.stopSimulation)
	stopButton.Importance = widget.LowImportance

	speedLabel := widget.NewLabel(speedString(3600).String())
	speedComponentsWidth := speedLabel.MinSize().Width
	speedLabel.SetText(speedString(sim.Speed()).String())

	speedSlider := widget.NewSlider(1, 3600)
	speedSlider.Step = 1
	speedSlider.Value = sim.Speed()
	speedSlider.OnChanged = controller.changeSpeedSimulation

	controller.startButton = startButton
	controller.stopButton = stopButton
	controller.speedLabel = speedLabel
	controller.speedSlider = speedSlider

	stickyRect := canvas.NewRectangle(color.Transparent)
	stickyRect.SetMinSize(fyne.NewSize(speedComponentsWidth, 0))
	return container.NewHBox(
		container.NewPadded(stopButton),
		container.NewVBox(stickyRect, layout.NewSpacer(), speedLabel, layout.NewSpacer(), stickyRect),
		container.NewVBox(stickyRect, layout.NewSpacer(), speedSlider, layout.NewSpacer(), stickyRect),
		container.NewPadded(startButton),
	)
}

func buildCityForm(city api.City, window fyne.Window) fyne.CanvasObject {
	nameItem := widget.NewFormItem("Name", widget.NewLabel(city.Name()))
	pos := city.Position().ToPos32()
	positionItem := widget.NewFormItem(
		"Position",
		widget.NewLabel(fmt.Sprintf("X: %d, Y: %d", int(pos.X), int(pos.Y))),
	)
	colorItem := widget.NewFormItem(
		"Color",
		buildColorChooser(city.Color, city.SetColor, window),
	)
	processingItem := widget.NewFormItem(
		"Processing Time",
		buildDurationSlider(city.ProcessingTime, city.SetProcessingTime),
	)
	generationItem := widget.NewFormItem(
		"Generation Time",
		buildDurationSlider(city.GenerationTime, city.SetGenerationTime),
	)

	stateItem := widget.NewFormItem(
		"State",
		buildPlayControlBar(city.Running, city.Start, city.Stop),
	)

	return widget.NewForm(nameItem, positionItem, colorItem, processingItem, generationItem, stateItem)
}

func buildVehicleForm(vehicle api.Vehicle, window fyne.Window) (fyne.CanvasObject, func()) {
	plateItem := widget.NewFormItem("Plate", widget.NewLabel(vehicle.Plate()))
	colorItem := widget.NewFormItem("Color", buildColorChooser(vehicle.Color, vehicle.SetColor, window))
	speedItem := widget.NewFormItem("Speed", buildSpeedSlider(vehicle.PreferredSpeed, vehicle.SetPreferredSpeed))
	bar := widget.NewProgressBar()
	progressItem := widget.NewFormItem("Progress", bar)
	stopCh := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second / 60)
		for {
			select {
			case <-stopCh:
				ticker.Stop()
				return
			case <-ticker.C:
				bar.Value = vehicle.Progress()
				bar.Refresh()
			}
		}
	}()

	return widget.NewForm(plateItem, colorItem, speedItem, progressItem), func() {
		stopCh <- struct{}{}
		close(stopCh)
	}
}

func buildPlayControlBar(running func() bool, play, stop func()) fyne.CanvasObject {
	playBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil)
	playBtn.Importance = widget.LowImportance
	stopBtn := widget.NewButtonWithIcon("", theme.MediaStopIcon(), nil)
	stopBtn.Importance = widget.LowImportance
	enableFn := func() {
		if running() {
			stopBtn.Enable()
			playBtn.Disable()
		} else {
			stopBtn.Disable()
			playBtn.Enable()
		}
	}
	playBtn.OnTapped = func() {
		play()
		enableFn()
	}
	stopBtn.OnTapped = func() {
		stop()
		enableFn()
	}
	enableFn()
	return container.NewHBox(stopBtn, playBtn)
}

func buildColorChooser(current func() color.Color, set func(color.Color), window fyne.Window) fyne.CanvasObject {
	rect := canvas.NewRectangle(current())

	btn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		dialog.ShowColorPicker("Choose Color", "", func(c color.Color) {
			set(c)
			rect.FillColor = c
			rect.Refresh()
		}, window)
	})

	return container.NewGridWrap(btn.MinSize(), rect, btn)
}

func buildDurationSlider(current func() time.Duration, set func(duration time.Duration)) fyne.CanvasObject {
	label := widget.NewLabel(current().String())
	label.Alignment = fyne.TextAlignCenter
	slider := widget.NewSlider(float64(time.Second/10), float64(time.Hour))
	slider.SetValue(float64(current()))
	slider.Step = float64(time.Second)
	slider.OnChanged = func(f float64) {
		d := time.Duration(f)
		set(d)
		label.SetText(d.String())
	}
	return container.NewVBox(label, slider)
}

func buildSpeedSlider(current func() float64, set func(value float64)) fyne.CanvasObject {
	format := func(value float64) string { return fmt.Sprintf("%dkm/h", int(value)) }
	label := widget.NewLabel(format(current()))
	label.Alignment = fyne.TextAlignCenter
	slider := widget.NewSlider(60, 600)
	slider.Step = 5
	slider.SetValue(current())
	slider.OnChanged = func(f float64) {
		set(f)
		label.SetText(format(f))
	}
	return container.NewVBox(label, slider)
}

func buildMenu(controller *simulationPlayController) *fyne.MainMenu {
	return fyne.NewMainMenu(buildSimulationMenu(controller))
}

func buildSimulationMenu(controller *simulationPlayController) *fyne.Menu {
	start := fyne.NewMenuItem("Start", controller.startSimulation)
	start.Icon = theme.MediaPlayIcon()
	stop := fyne.NewMenuItem("Stop", controller.stopSimulation)
	stop.Icon = theme.MediaStopIcon()
	speed := fyne.NewMenuItem("Speed", func() {
		stickyRect := canvas.NewRectangle(color.Transparent)
		stickyRect.SetMinSize(fyne.NewSize(controller.speedLabel.MinSize().Height, 0))
		cnt := container.NewVBox(
			stickyRect,
			container.NewWithoutLayout(controller.speedLabel),
			container.NewWithoutLayout(controller.speedSlider),
			stickyRect, stickyRect, stickyRect, stickyRect, stickyRect,
		)
		dialog.ShowCustom("Change Simulation Speed", "Cancel", cnt, controller.mainWindow)
	})
	speed.Icon = theme.MediaFastForwardIcon()

	controller.startMenuButton = start
	controller.stopMenuButton = stop

	return fyne.NewMenu("Simulation", start, stop, speed)
}

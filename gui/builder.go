package gui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com.bisoncorp.autostrade/game"
	api "github.com.bisoncorp.autostrade/gameapi"
	"github.com.bisoncorp.autostrade/gui/controller"
	gamewid "github.com.bisoncorp.autostrade/gui/widget"
	"github.com.bisoncorp.autostrade/sampledata"
	"image/color"
	"io"
	"time"
)

func buildSimulationUi(sim api.Simulation, window fyne.Window, application *Application) (fyne.CanvasObject, *fyne.MainMenu) {
	simulationRunnableController := controller.NewRunnableController(sim)
	simulationSpeedableController := controller.NewSpeedableController(sim)
	hintController, hintObject := controller.NewHintController()

	leftCnt, addCity := buildCityPropertiesContainer(window)
	rightCnt, addVehicle := buildVehiclesPropertiesContainer(window)

	mapObject, mapWidget := buildMap()

	mapWidget.OnCityTapped = func(data api.CityData) {
		c := sim.City(data.Name)
		addCity(c)
	}
	mapWidget.OnVehicleTapped = func(data api.VehicleData) {
		v := sim.Vehicle(data.Plate)
		addVehicle(v)
	}
	mapWidget.OnDataRequired = func() api.SimulationData {
		return sim.PackData()
	}

	return container.NewBorder(
		buildToolbar(sim, mapWidget, window, hintController),
		container.NewBorder(nil, nil, nil, hintObject, buildSimulationControlBar(simulationRunnableController, simulationSpeedableController)),
		leftCnt, rightCnt,
		mapObject,
	), buildMenu(sim, simulationRunnableController, simulationSpeedableController, window, application)
}

func buildMap() (fyne.CanvasObject, *gamewid.Map) {
	mapWidget := gamewid.NewMap()

	background := canvas.NewImageFromImage(sampledata.ItalyMap())
	size := background.Image.Bounds().Size()
	fSize := fyne.NewSize(float32(size.X), float32(size.Y))
	ratio := fSize.Width / fSize.Height
	background.SetMinSize(fyne.NewSize(600*ratio, 600))
	background.Resize(background.MinSize())

	cnt := container.NewMax(container.NewWithoutLayout(background), mapWidget)
	scroll := container.NewScroll(cnt)
	scroll.SetMinSize(fyne.NewSize(600*ratio, 600))
	return scroll, mapWidget
}

func buildCityPropertiesContainer(window fyne.Window) (obj fyne.CanvasObject, add func(api.City)) {
	cities := make(map[api.City]int)
	accordion := widget.NewAccordion()
	addCity := func(city api.City) {
		if index, exist := cities[city]; exist {
			accordion.Open(index)
			return
		}

		title := fmt.Sprintf("City Property [%s]", city.Name())
		item := widget.NewAccordionItem(title, nil)
		content := buildCityProperty(city, window)
		closeBtn := widget.NewButtonWithIcon("Close", theme.CancelIcon(), func() {
			index := cities[city]
			for k, v := range cities {
				if v > index {
					cities[k]--
				}
			}
			delete(cities, city)
			accordion.Remove(item)
		})
		closeBtn.Importance = widget.LowImportance
		item.Detail = container.NewVBox(content, closeBtn)

		index := len(accordion.Items)
		cities[city] = index
		accordion.CloseAll()
		accordion.Append(item)
		accordion.Open(index)
	}
	return accordion, addCity
}
func buildVehiclesPropertiesContainer(window fyne.Window) (obj fyne.CanvasObject, add func(api.Vehicle)) {
	vehicles := make(map[api.Vehicle]int)
	accordion := widget.NewAccordion()
	addVehicle := func(vehicle api.Vehicle) {
		if index, exist := vehicles[vehicle]; exist {
			accordion.Open(index)
			return
		}

		title := fmt.Sprintf("Vehicle Property [%s]", vehicle.Plate())
		item := widget.NewAccordionItem(title, nil)
		content, closeView := buildVehicleProperty(vehicle, window)
		closeBtn := widget.NewButtonWithIcon("Close", theme.CancelIcon(), func() {
			index := vehicles[vehicle]
			for k, v := range vehicles {
				if v > index {
					vehicles[k]--
				}
			}
			delete(vehicles, vehicle)
			accordion.Remove(item)
			closeView()
		})
		closeBtn.Importance = widget.LowImportance
		item.Detail = container.NewVBox(content, closeBtn)

		index := len(accordion.Items)
		vehicles[vehicle] = index
		accordion.CloseAll()
		accordion.Append(item)
		accordion.Open(index)
	}
	return accordion, addVehicle
}

func buildSimulationControlBar(rc *controller.RunnableController, sc *controller.SpeedableController) fyne.CanvasObject {
	startButton, stopButton := buildRunnableControlBar(rc)
	speedBar := buildSpeedableControlBar(sc)
	return container.NewBorder(nil, nil, container.NewHBox(stopButton, startButton), nil, speedBar)
}
func buildSpeedableControlBar(ctrl *controller.SpeedableController) fyne.CanvasObject {
	label := widget.NewLabel(speedString(3600).String())
	slider := widget.NewSlider(1, 3600)
	size := label.MinSize()
	label.SetText(speedString(ctrl.Speed()).String())
	slider.SetValue(ctrl.Speed())
	slider.OnChanged = func(f float64) {
		ctrl.SetSpeed(f)
	}
	ctrl.AddCallback(func(f float64) {
		label.SetText(speedString(f).String())
		slider.Value = f
		slider.Refresh()
	})
	return container.NewGridWrap(size, label, slider)
}
func buildRunnableControlBar(ctrl *controller.RunnableController) (start fyne.CanvasObject, stop fyne.CanvasObject) {
	playBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil)
	playBtn.Importance = widget.LowImportance
	stopBtn := widget.NewButtonWithIcon("", theme.MediaStopIcon(), nil)
	stopBtn.Importance = widget.LowImportance

	enableFn := func(en bool) {
		if en {
			stopBtn.Enable()
			playBtn.Disable()
		} else {
			playBtn.Enable()
			stopBtn.Disable()
		}
	}

	playBtn.OnTapped = func() {
		playBtn.Disable()
		ctrl.Start()
	}
	stopBtn.OnTapped = func() {
		stopBtn.Disable()
		ctrl.Stop()
	}

	enableFn(ctrl.Running())
	ctrl.AddCallback(func(eventType controller.RunnableEventType) { enableFn(eventType == controller.Started) })
	return playBtn, stopBtn
}

func buildCityProperty(city api.City, window fyne.Window) fyne.CanvasObject {
	nameItem := widget.NewFormItem("Name", widget.NewLabel(city.Name()))
	pos := city.Position().ToPos32()
	positionItem := widget.NewFormItem(
		"Position",
		widget.NewLabel(fmt.Sprintf("X: %d, Y: %d", int(pos.X), int(pos.Y))),
	)
	colorItem := widget.NewFormItem(
		"Color",
		buildColorChooser(controller.NewColorableController(city), window),
	)
	processingItem := widget.NewFormItem(
		"Processing Time",
		buildDurationSlider(city.ProcessingTime, city.SetProcessingTime),
	)
	generationItem := widget.NewFormItem(
		"Generation Time",
		buildDurationSlider(city.GenerationTime, city.SetGenerationTime),
	)

	start, stop := buildRunnableControlBar(controller.NewRunnableController(city))
	stateItem := widget.NewFormItem(
		"State",
		container.NewHBox(stop, start),
	)

	return widget.NewForm(nameItem, positionItem, colorItem, processingItem, generationItem, stateItem)
}
func buildVehicleProperty(vehicle api.Vehicle, window fyne.Window) (obj fyne.CanvasObject, clear func()) {
	plateItem := widget.NewFormItem("Plate", widget.NewLabel(vehicle.Plate()))
	colorItem := widget.NewFormItem("Color", buildColorChooser(controller.NewColorableController(vehicle), window))
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

func buildColorChooser(ctrl *controller.ColorableController, window fyne.Window) fyne.CanvasObject {
	rect := canvas.NewRectangle(ctrl.Color())
	ctrl.AddCallback(func(c color.Color) {
		rect.FillColor = c
		rect.Refresh()
	})

	btn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		picker := dialog.NewColorPicker("Choose Color", "", func(c color.Color) {
			ctrl.SetColor(c)
		}, window)
		picker.Advanced = true
		picker.Show()
	})

	return container.NewGridWrap(btn.MinSize(), rect, btn)
}
func buildDurationSlider(get func() time.Duration, set func(duration time.Duration)) fyne.CanvasObject {
	label := widget.NewLabel(get().String())
	label.Alignment = fyne.TextAlignCenter
	slider := widget.NewSlider(float64(time.Second/10), float64(time.Hour))
	slider.SetValue(float64(get()))
	slider.Step = float64(time.Second)
	slider.OnChanged = func(f float64) {
		d := time.Duration(f)
		set(d)
		label.SetText(d.String())
	}
	return container.NewVBox(label, slider)
}
func buildSpeedSlider(get func() float64, set func(value float64)) fyne.CanvasObject {
	format := func(value float64) string { return fmt.Sprintf("%dkm/h", int(value)) }
	label := widget.NewLabel(format(get()))
	label.Alignment = fyne.TextAlignCenter
	slider := widget.NewSlider(60, 600)
	slider.Step = 5
	slider.SetValue(get())
	slider.OnChanged = func(f float64) {
		set(f)
		label.SetText(format(f))
	}
	return container.NewVBox(label, slider)
}

func buildToolbar(sim api.Simulation, mapWidget *gamewid.Map, window fyne.Window, hintController *controller.HintController) fyne.CanvasObject {
	addCityBtn := widget.NewButton("Add City", nil)
	addCityBtn.Importance = widget.LowImportance

	remCityBtn := widget.NewButton("Remove City", nil)
	remCityBtn.Importance = widget.LowImportance

	moveCityBtn := widget.NewButton("Move City", nil)
	moveCityBtn.Importance = widget.LowImportance

	addRoadBtn := widget.NewButton("Add Road", nil)
	addRoadBtn.Importance = widget.LowImportance

	disableAll := func() {
		addRoadBtn.Disable()
		addCityBtn.Disable()
		remCityBtn.Disable()
		moveCityBtn.Disable()
	}

	enableAll := func() {
		addRoadBtn.Enable()
		addCityBtn.Enable()
		remCityBtn.Enable()
		moveCityBtn.Enable()
	}

	addCityBtn.OnTapped = func() {
		go func() {
			disableAll()
			<-actionAddCity(sim, mapWidget, window, hintController)
			enableAll()
		}()
	}
	remCityBtn.OnTapped = func() {
		go func() {
			disableAll()
			<-actionRemCity(sim, mapWidget, hintController)
			enableAll()
		}()
	}
	moveCityBtn.OnTapped = func() {
		go func() {
			disableAll()
			<-actionMoveCity(sim, mapWidget, hintController)
			enableAll()
		}()
	}
	addRoadBtn.OnTapped = func() {
		go func() {
			disableAll()
			<-actionAddRoad(sim, mapWidget, window, hintController)
			enableAll()
		}()
	}
	return container.NewHBox(addCityBtn, remCityBtn, moveCityBtn, widget.NewSeparator(), addRoadBtn)
}

func buildMenu(sim api.Simulation, rc *controller.RunnableController, sc *controller.SpeedableController, window fyne.Window, application *Application) *fyne.MainMenu {
	return fyne.NewMainMenu(buildFileMenu(sim, window, application), buildSimulationMenu(rc, sc, window))
}
func buildSimulationMenu(rc *controller.RunnableController, sc *controller.SpeedableController, window fyne.Window) *fyne.Menu {
	start := fyne.NewMenuItem("Start", rc.Start)
	start.Icon = theme.MediaPlayIcon()
	stop := fyne.NewMenuItem("Stop", rc.Stop)
	stop.Icon = theme.MediaStopIcon()

	rc.AddCallback(func(eventType controller.RunnableEventType) {
		val := eventType == controller.Started
		start.Disabled = val
		stop.Disabled = !val
	})

	speedControl := buildSpeedableControlBar(sc)
	speed := fyne.NewMenuItem("Speed", func() {
		dialog.ShowCustom("Change Simulation Speed", "Cancel", speedControl, window)
	})
	speed.Icon = theme.MediaFastForwardIcon()

	return fyne.NewMenu("Simulation", start, stop, speed)
}
func buildFileMenu(sim api.Simulation, window fyne.Window, application *Application) *fyne.Menu {
	var writer io.Writer
	saveWithName := func() {
		dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
			if file == nil {
				return
			}
			writer = file
			data := sim.PackData()
			data.Vehicles = make([]struct {
				api.VehicleData
				RoadIndex int
			}, 0)
			api.WriteSimulationData(data, writer)
		}, window)
	}
	saveWithNameItem := fyne.NewMenuItem("Save With Name", saveWithName)
	saveWithNameItem.Icon = theme.DocumentSaveIcon()

	saveItem := fyne.NewMenuItem("Save", func() {
		if writer == nil {
			saveWithName()
			return
		}
		data := sim.PackData()
		data.Vehicles = make([]struct {
			api.VehicleData
			RoadIndex int
		}, 0)
		api.WriteSimulationData(data, writer)

	})
	saveItem.Icon = theme.DocumentSaveIcon()

	openItem := fyne.NewMenuItem("Open", func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if file == nil {
				return
			}
			data := api.ReadSimulationData(file)
			ns := game.NewFromData(data)
			application.NewWindow(ns)
		}, window)
	})
	openItem.Icon = theme.FileIcon()

	newItem := fyne.NewMenuItem("New", func() {
		application.NewWindow(game.New())
	})
	newItem.Icon = theme.ContentAddIcon()

	return fyne.NewMenu("File", saveItem, saveWithNameItem, openItem, newItem)
}

func showCityForm(sim api.Simulation, window fyne.Window) <-chan api.CityData {
	nameEntry := widget.NewEntry()
	nameEntry.Validator = func(s string) error {
		city := sim.City(s)
		if city != nil {
			return errors.New("city already exist")
		}
		return nil
	}
	nameEntry.PlaceHolder = sampledata.RandomCityName()
	nameItem := widget.NewFormItem("Name", nameEntry)

	colorBuffer := controller.NewColorableController(controller.NewColorableBuffer(randomColor()))
	colorItem := widget.NewFormItem(
		"Color",
		buildColorChooser(colorBuffer, window),
	)

	processingDuration := time.Millisecond * 100
	processingItem := widget.NewFormItem(
		"Processing Time",
		buildDurationSlider(func() time.Duration {
			return processingDuration
		}, func(duration time.Duration) {
			processingDuration = duration
		}),
	)

	generationDuration := time.Millisecond * 500
	generationItem := widget.NewFormItem(
		"Generation Time",
		buildDurationSlider(func() time.Duration {
			return generationDuration
		}, func(duration time.Duration) {
			generationDuration = duration
		}),
	)

	items := []*widget.FormItem{nameItem, colorItem, processingItem, generationItem}
	ch := make(chan api.CityData, 1)
	dialog.ShowForm("New City", "Choose Position", "Cancel", items, func(confirmed bool) {
		if !confirmed {
			close(ch)
			return
		}
		name := nameEntry.PlaceHolder
		if nameEntry.Text != "" {
			name = nameEntry.Text
		}
		ch <- api.CityData{
			Name:           name,
			Color:          colorToRgba(colorBuffer.Color()),
			GenerationTime: generationDuration,
			ProcessingTime: processingDuration,
		}
		close(ch)
	}, window)

	return ch
}
func showRoadForm(sim api.Simulation, window fyne.Window) <-chan struct {
	data   api.RoadData
	oneWay bool
} {
	maxSpeed := float64(130)
	slider := buildSpeedSlider(func() float64 {
		return maxSpeed
	}, func(value float64) {
		maxSpeed = value
	})
	maxSpeedItem := widget.NewFormItem("Max Speed", slider)

	oneWay := false
	check := widget.NewCheck("One Way?", func(b bool) {
		oneWay = b
	})
	oneWayItem := widget.NewFormItem("", check)

	items := []*widget.FormItem{maxSpeedItem, oneWayItem}
	ch := make(chan struct {
		data   api.RoadData
		oneWay bool
	}, 1)
	dialog.ShowForm("New Road", "Next", "Cancel", items, func(confirmed bool) {
		if !confirmed {
			close(ch)
			return
		}
		ch <- struct {
			data   api.RoadData
			oneWay bool
		}{data: api.RoadData{MaxSpeed: maxSpeed}, oneWay: oneWay}
		close(ch)
	}, window)

	return ch
}

func actionAddCity(sim api.Simulation, mapWidget *gamewid.Map, window fyne.Window, hintController *controller.HintController) (done chan struct{}) {
	done = make(chan struct{}, 1)
	go func() {
		dataCh := showCityForm(sim, window)
		data, ok := <-dataCh
		if !ok {
			return
		}

		posCh := make(chan fyne.Position)
		defer close(posCh)
		mapWidget.OnTapped = func(event *fyne.PointEvent) {
			posCh <- event.Position
		}
		defer func() { mapWidget.OnTapped = nil }()

		hintController.SetHint("Choose position")
		pos := <-posCh
		hintController.Clear()
		data.Pos = api.Position{X: float64(pos.X), Y: float64(pos.Y)}
		if city := sim.AddCity(data); sim.Running() {
			city.Start()
		}
		done <- struct{}{}
	}()
	return
}
func actionAddRoad(sim api.Simulation, mapWidget *gamewid.Map, window fyne.Window, hintController *controller.HintController) (done chan struct{}) {
	done = make(chan struct{}, 1)
	go func() {
		dataCh := showRoadForm(sim, window)
		data, ok := <-dataCh
		if !ok {
			return
		}

		cityCh := make(chan api.City)
		defer close(cityCh)
		oldFn := mapWidget.OnCityTapped
		mapWidget.OnCityTapped = func(data api.CityData) {
			city := sim.City(data.Name)
			if city == nil {
				panic("city is nil, unexpected")
			}
			cityCh <- city
		}
		defer func() { mapWidget.OnCityTapped = oldFn }()

		hintController.SetHint("Select first city")
		city1 := <-cityCh
		hintController.SetHint("Select second city")
		city2 := <-cityCh
		hintController.Clear()
		if data.oneWay {
			if road := sim.AddOneWayRoad(city1, city2, data.data); sim.Running() {
				if road != nil {
					road.Start()
				}
			}
		} else {
			if atob, btoa := sim.AddRoad(city1, city2, data.data); sim.Running() {
				if atob != nil {
					atob.Start()
				}
				if btoa != nil {
					btoa.Start()
				}
			}
		}
		done <- struct{}{}
	}()
	return
}
func actionRemCity(sim api.Simulation, mapWidget *gamewid.Map, hintController *controller.HintController) (done chan struct{}) {
	done = make(chan struct{}, 1)
	go func() {
		ch := make(chan string)
		oldFn := mapWidget.OnCityTapped
		mapWidget.OnCityTapped = func(data api.CityData) {
			ch <- data.Name
		}
		defer func() { mapWidget.OnCityTapped = oldFn }()

		hintController.SetHint("Select a city")
		sim.RemoveCity(sim.City(<-ch))
		hintController.Clear()
		done <- struct{}{}
	}()
	return
}
func actionMoveCity(sim api.Simulation, mapWidget *gamewid.Map, hintController *controller.HintController) (done chan struct{}) {
	done = make(chan struct{}, 1)
	go func() {
		chName, chPos := make(chan string), make(chan fyne.Position)
		defer close(chName)
		defer close(chPos)
		oldFn := mapWidget.OnCityTapped
		mapWidget.OnCityTapped = func(data api.CityData) {
			chName <- data.Name
		}
		defer func() { mapWidget.OnCityTapped = oldFn }()
		mapWidget.OnTapped = func(event *fyne.PointEvent) {
			chPos <- event.Position
		}
		defer func() { mapWidget.OnTapped = nil }()

		hintController.SetHint("Select a city")
		city := sim.City(<-chName)
		hintController.SetHint("Choose new position")
		pos := <-chPos
		hintController.Clear()

		x, y := pos.Components()
		city.SetPosition(api.Position{X: float64(x), Y: float64(y)})

		done <- struct{}{}
	}()
	return
}

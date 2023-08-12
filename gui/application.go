package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com.bisoncorp.autostrade/game"
	api "github.com.bisoncorp.autostrade/gameapi"
	"os"
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

func (a *Application) NewWindow(path string) {
	wind := a.appl.NewWindow("Simulation")
	wind.Resize(fyne.NewSize(600, 600))
	wind.CenterOnScreen()
	var sim api.Simulation
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		data := api.ReadSimulationData(file)
		err = file.Close()
		if err != nil {
			panic(err)
		}
		sim = game.NewFromData(data)
	} else {
		sim = game.New()
	}
	ui, menu := buildSimulationUi(sim, wind)
	items := menu.Items
	fileMenu := buildFileMenu(sim, wind, path, a)
	menu.Items = append([]*fyne.Menu{}, fileMenu)
	menu.Items = append(menu.Items, items...)

	wind.SetContent(ui)
	wind.SetMainMenu(menu)
	wind.Show()
}

func saveWithName(simulation api.Simulation, window fyne.Window) <-chan string {
	ch := make(chan string)
	go dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		defer close(ch)
		if writer == nil {
			return
		}
		data := simulation.PackData()
		data.Vehicles = make([]struct {
			api.VehicleData
			RoadIndex int
		}, 0)
		api.WriteSimulationData(data, writer)
		ch <- writer.URI().Path()
	}, window)
	return ch
}

func buildFileMenu(simulation api.Simulation, window fyne.Window, path string, application *Application) *fyne.Menu {
	saveWithNameItem := fyne.NewMenuItem("Save With Name", func() {
		go func() {
			p, ok := <-saveWithName(simulation, window)
			if ok {
				path = p
			}
		}()
	})
	saveWithNameItem.Icon = theme.DocumentSaveIcon()

	saveItem := fyne.NewMenuItem("Save", func() {
		if path == "" {
			go func() {
				p, ok := <-saveWithName(simulation, window)
				if ok {
					path = p
				}
			}()
			return
		}
		data := simulation.PackData()
		data.Vehicles = make([]struct {
			api.VehicleData
			RoadIndex int
		}, 0)
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		api.WriteSimulationData(data, file)
		err = file.Close()
		if err != nil {
			panic(err)
		}

	})
	saveItem.Icon = theme.DocumentSaveIcon()

	openItem := fyne.NewMenuItem("Open", func() {
		dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
			if f == nil {
				return
			}
			application.NewWindow(f.URI().Path())
		}, window)
	})
	openItem.Icon = theme.FileIcon()

	newItem := fyne.NewMenuItem("New", func() {
		application.NewWindow("")
	})
	newItem.Icon = theme.ContentAddIcon()

	return fyne.NewMenu("File", saveItem, saveWithNameItem, openItem, newItem)
}

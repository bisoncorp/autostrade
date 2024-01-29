package main

import (
	"github.com/bisoncorp/autostrade/game"
	api "github.com/bisoncorp/autostrade/gameapi"
	"github.com/bisoncorp/autostrade/gui"
	"log"
	"os"
)

func main() {
	appl := gui.NewApplication()
	filesPath := os.Args[1:]
	if len(filesPath) > 0 {
		for _, path := range filesPath {
			file, err := os.Open(path)
			if err != nil {
				log.Println(err)
				continue
			}
			appl.NewWindow(game.NewFromData(api.ReadSimulationData(file)))
			_ = file.Close()
		}
	} else {
		appl.NewWindow(game.New())
	}
	appl.Run()
}

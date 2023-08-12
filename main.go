package main

import (
	"github.com.bisoncorp.autostrade/gui"
	"os"
)

func main() {
	appl := gui.NewApplication()
	filesPath := os.Args[1:]
	if len(filesPath) > 0 {
		for _, path := range filesPath {
			appl.NewWindow(path)
		}
	} else {
		appl.NewWindow("")
	}
	appl.Run()
}

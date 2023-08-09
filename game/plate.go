package game

import (
	api "github.com.bisoncorp.autostrade/gameapi"
)

var (
	plateCh chan string
)

func init() {
	plateCh = make(chan string, 16)
	go plateGen()
}

func plateGen() {
	plate := api.FirstPlate
	for {
		plateCh <- plate.String()
		plate = plate.Next()
	}
}

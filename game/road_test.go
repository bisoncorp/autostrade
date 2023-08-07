package game

import (
	api "github.com.bisoncorp.autostrade/gameapi"
	"testing"
	"time"
)

func Test_newRoad(t *testing.T) {
	r := newRoad(api.RoadData{}, nil)
	r.Start()
	time.Sleep(2 * time.Second)
	r.Stop()
	time.Sleep(2 * time.Second)
}

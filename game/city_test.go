package game

import (
	api "github.com.bisoncorp.autostrade/gameapi"
	"testing"
	"time"
)

func Test_city_Tick(t *testing.T) {
	c := newCity(api.CityData{})
	c.Start()
	time.Sleep(2 * time.Second)
	c.Stop()
	time.Sleep(2 * time.Second)
}

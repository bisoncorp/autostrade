package game

import "testing"

func Test_plateGen(t *testing.T) {
	for i := 0; i < 1<<10; i++ {
		t.Log(<-plateCh)
	}

}

package game

import "fmt"

var (
	plateCh chan string
)

func init() {
	plateCh = make(chan string, 16)
	go plateGen()
}

func plateGen() {
	a, b, c, d, n := 'A', 'A', 'A', 'A', 0
	increment := func() {
		n++
		if n == 1000 {
			n = 0
			d++
			if d == 'Z'+1 {
				d = 'A'
				c++
				if c == 'Z'+1 {
					c = 'A'
					b++
					if b == 'Z'+1 {
						b = 'A'
						a++
						if a == 'Z'+1 {
							a = 'A'
						}
					}
				}
			}
		}
	}
	for {
		plate := fmt.Sprintf("%c%c%03d%c%c", a, b, n, c, d)
		increment()
		plateCh <- plate
	}
}

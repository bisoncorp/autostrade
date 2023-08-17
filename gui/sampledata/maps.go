package sampledata

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
)

//go:embed map_of_italy.png
var mapOfItaly []byte

func ItalyMap() image.Image {
	decode, _, err := image.Decode(bytes.NewReader(mapOfItaly))
	if err != nil {
		panic(err)
	}
	return decode
}

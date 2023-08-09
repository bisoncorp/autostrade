package controller

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HintController struct {
	label *widget.Label
}

func NewHintController() (*HintController, fyne.CanvasObject) {
	h := &HintController{label: widget.NewLabel("")}
	h.Clear()
	return h, h.label
}

func (h *HintController) SetHint(hint string) {
	h.label.Show()
	h.label.SetText(fmt.Sprintf("Hint: %s", hint))
}

func (h *HintController) Clear() {
	h.label.Hide()
}

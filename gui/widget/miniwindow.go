package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MiniWindow struct {
	widget.BaseWidget

	Title   string
	Content fyne.CanvasObject
	onClose func(*MiniWindow)
}

func NewMiniWindow(title string, content fyne.CanvasObject, onClose func(*MiniWindow)) *MiniWindow {
	m := &MiniWindow{Title: title, Content: content, onClose: onClose}
	m.ExtendBaseWidget(m)
	m.Refresh()
	return m
}

func (m *MiniWindow) SetTitle(title string) {
	m.ExtendBaseWidget(m)
	m.Title = title
	m.Refresh()
}

func (m *MiniWindow) SetContent(content fyne.CanvasObject) {
	m.ExtendBaseWidget(m)
	m.Content = content
	m.Refresh()
}

func (m *MiniWindow) CreateRenderer() fyne.WidgetRenderer {
	titleLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), nil)
	showContentButton := widget.NewButtonWithIcon("", theme.MenuDropUpIcon(), nil)

	closeButton.Importance = widget.LowImportance
	showContentButton.Importance = widget.LowImportance

	closeButton.OnTapped = func() {
		if m.onClose != nil {
			m.onClose(m)
		}
	}
	showContentButton.OnTapped = func() {
		if c := m.Content; c != nil {
			if c.Visible() {
				c.Hide()
			} else {
				c.Show()
			}
			m.Refresh()
		}
	}

	contentCnt := container.NewPadded()
	window := container.NewBorder(
		container.NewHBox(titleLabel, layout.NewSpacer(), showContentButton, closeButton),
		nil, nil, nil, contentCnt)
	return &miniWindowRenderer{
		widget:            m,
		titleLabel:        titleLabel,
		showContentButton: showContentButton,
		contentCnt:        contentCnt,
		window:            window,
	}
}

type miniWindowRenderer struct {
	widget            *MiniWindow
	titleLabel        *widget.Label
	showContentButton *widget.Button
	contentCnt        *fyne.Container
	window            *fyne.Container
}

func (m *miniWindowRenderer) Destroy()                     {}
func (m *miniWindowRenderer) Layout(size fyne.Size)        { m.window.Resize(size) }
func (m *miniWindowRenderer) MinSize() fyne.Size           { return m.window.MinSize() }
func (m *miniWindowRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{m.window} }
func (m *miniWindowRenderer) Refresh() {
	m.titleLabel.SetText(m.widget.Title)
	if c := m.widget.Content; c != nil {
		m.showContentButton.Show()
		m.contentCnt.Objects = []fyne.CanvasObject{c}
		m.contentCnt.Refresh()
		if c.Visible() {
			m.showContentButton.SetIcon(theme.MenuDropUpIcon())
		} else {
			m.showContentButton.SetIcon(theme.MenuDropDownIcon())
		}
	} else {
		m.showContentButton.Hide()
	}
}

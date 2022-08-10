package ui

import "fyne.io/fyne/v2"

type UiPage interface {
	renderUi() fyne.CanvasObject
	pageSize() fyne.Size
	isNewWindow() bool
	pageName() string
	setMainFrame(*MainFrame)
	setWindow(w fyne.Window)
}

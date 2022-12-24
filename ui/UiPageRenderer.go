package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func NewMainFrame(frameName string) MainFrame {
	mf := MainFrame{}

	mf.app = app.New()
	mf.masterWindow = mf.app.NewWindow(frameName)

	return mf
}

type MainFrame struct {
	masterWindow fyne.Window
	app          fyne.App
}

func (mf *MainFrame) RenderUiPage(p UiPage) {
	p.setMainFrame(mf)

	var pageWindow fyne.Window
	page := p.renderUi()

	if p.isNewWindow() {
		pageWindow = mf.app.NewWindow(p.pageName())
		pageWindow.Show()
	} else {
		pageWindow = mf.masterWindow
	}

	pageSize := p.pageSize()
	if pageSize != (fyne.Size{}) {
		pageWindow.Resize(pageSize)
	}

	pageName := p.pageName()
	if pageName != "" {
		pageWindow.SetTitle(pageName)
	}

	p.setWindow(pageWindow)
	pageWindow.SetContent(page)
}

func (mf *MainFrame) ShowAndRun() {
	mf.masterWindow.SetMaster()
	mf.masterWindow.ShowAndRun()
}
